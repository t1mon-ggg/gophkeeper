package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mgutz/ansi"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"github.com/t1mon-ggg/gophkeeper/pkg/client/storage"
	"github.com/t1mon-ggg/gophkeeper/pkg/client/storage/secrets"
)

func (c *CLI) save() {
	if c.config.Mode != "standalone" {
		fmt.Println("todo save to remote")
	}
	buf, err := c.storage.Save()
	if err != nil {
		if errors.Is(err, storage.ErrHashValid) {
			c.log().Info(nil, "skip saving")
			return
		}
		c.log().Error(err, "export failed")
		return
	}
	f, err := os.OpenFile(c.config.Storage, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		c.log().Error(err, "open file failed")
		return
	}
	defer f.Close()
	msg, err := c.crypto.EncryptWithKeys(buf)
	if err != nil {
		c.log().Error(err, "encrypting failed")
		return
	}
	_, err = f.Write(msg)
	if err != nil {
		c.log().Error(err, "write to file failed")
		return
	}
}

func (c *CLI) insert(in string) {
	args := strings.Split(in, " ")
	if len(args) < 3 {
		c.log().Error(nil, "not enought arguments")
		return
	}
	switch args[2] {
	case "otp":
		if len(args) < 15 {
			c.log().Error(nil, "not enought arguments")
			return
		}
		var secret, issuer, username, method, name, description string
		var codes []string
		name = args[1]
		for i, v := range args[1:] {
			switch v {
			case "method":
				secret = args[i+2]
			case "secret":
				secret = args[i+2]
			case "issuer":
				issuer = args[i+2]
			case "username":
				username = args[i+2]
			case "recoverycodes":
				for j, vv := range args[i+2:] {
					if vv == "description" {
						break
					}
					codes = append(codes, args[i+j+2])
				}
			case "description":
				description = strings.Join(args[i+2:], " ")
			}
		}
		if method == "" {
			method = "TOTP"
		}
		otp := secrets.NewOTP(method, issuer, secret, username, codes...)
		c.storage.InsertSecret(name, description, otp)
		return
	case "userpass":
		if len(args) < 8 {
			c.log().Error(nil, "not enought arguments")
			return
		}
		var username, password, name, description string
		name = args[1]
		for i, v := range args[1:] {
			switch v {
			case "username":
				username = args[i+2]
			case "password":
				password = args[i+2]
			case "description":
				description = strings.Join(args[i+2:], " ")
			}
		}
		up := secrets.NewUserPass(username, password)
		c.storage.InsertSecret(name, description, up)
		return
	case "creditcard":
		if len(args) != 13 {
			var number, holder, expire, name, description string
			var cvv int
			name = args[1]
			for i, v := range args[1:] {
				switch v {
				case "description":
					description = strings.Join(args[i+2:], " ")
				case "holder":
					holder = fmt.Sprintf("%s %s", args[i+2], args[i+3])
				case "expire":
					expire = args[i+2]
				case "number":
					number = args[i+2]
				case "cvv":
					c, err := strconv.Atoi(args[i+2])
					if err != nil {
						fmt.Println("credit card invalid")
						return
					}
					cvv = c
				}
			}
			s, err := secrets.NewCC(number, holder, expire, uint16(cvv))
			if err != nil {
				fmt.Println("credit card invalid")
				return
			}
			c.storage.InsertSecret(name, description, s)
		}
		return
	case "anytext":
		var name, text, description string
		var txt, dscr []string
		name = args[1]
		for i, v := range args[1:] {
			if v == "text" {
				for j, vv := range args[i+2:] {
					if vv == "description" {
						break
					}
					txt = append(txt, args[i+j+2])
				}
			}
			if v == "description" {
				for j, vv := range args[i+2:] {
					if vv == "text" {
						break
					}
					dscr = append(dscr, args[i+j+2])
				}
			}
		}
		text = strings.Join(txt, " ")
		description = strings.Join(dscr, " ")
		secret := secrets.NewText(text)
		c.storage.InsertSecret(name, description, secret)
		return
	case "anybinary":
		var name, path, description string
		name = args[1]
		for i, v := range args[1:] {
			if v == "path" {
				path = args[i+2]
			}
			if v == "description" {
				description = strings.Join(args[i+2:], " ")
			}
		}
		f, err := os.Open(path)
		if err != nil {
			c.log().Error(err, "secret file open failed")
			return
		}
		content, err := io.ReadAll(f)
		if err != nil {
			c.log().Error(err, "secret file read failed")
			return
		}
		bin := secrets.NewBinary(content)
		c.storage.InsertSecret(name, description, bin)
		return
	}
}

func (c *CLI) list() {
	lst := c.storage.ListSecrets()
	if len(lst) == 0 {
		fmt.Println("no secrets found")
		return
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Secret name", "Type", "Description"})
	i := 1
	for k, v := range lst {
		secret := c.storage.GetSecret(k)
		t.AppendRows([]table.Row{
			{i, k, secret.Scope(), v},
		})
		t.AppendSeparator()
		i++
	}
	t.SetStyle(table.StyleColoredBright)
	t.Render()
}

func (c *CLI) get(name string, opts ...string) {
	secret := c.storage.GetSecret(name)
	if secret == nil {
		fmt.Println("no such secret")
		return
	}
	scope := secret.Scope()
	switch scope {
	case "user-password":
		value := secret.Value().(*secrets.UserPass)
		fmt.Printf("Username: %s\nPassword: %s\n", value.Username, value.Password)
	case "anybinary":
		if len(opts) != 1 {
			c.log().Error(nil, "path to save binary not provided")
			return
		}
		f, err := os.OpenFile(opts[0], os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
		if err != nil {
			c.log().Error(err, "secret can not be saved")
			return
		}
		defer f.Close()
		value := secret.Value().(*secrets.AnyBinary)
		f.Write(value.Bytes)
	case "anytext":
		value := secret.Value().(*secrets.AnyText)
		fmt.Printf("Text:\n %s\n", value.Text)
	case "creditcard":
		value := secret.Value().(*secrets.CreditCard)
		t := time.Now()
		if t.After(value.Expire) {
			fmt.Println(ansi.Color("Credit card expired", "red+b"))
		}
		fmt.Printf("Number %s\nCardHolder: %s\nExpire: %s\tCVV: %v\n", value.Number, value.Holder, value.Expire.Format("01/06"), value.CVV)
	case "otp":
		value := secret.Value().(*secrets.OTP)
		switch value.Method {
		case "TOTP":
			code, err := totp.GenerateCodeCustom(value.Secret, time.Now(), totp.ValidateOpts{
				Period:    30,
				Skew:      1,
				Digits:    otp.DigitsSix,
				Algorithm: otp.AlgorithmSHA1,
			})
			if err != nil {
				c.log().Error(err, "secret can not be displayed")
			}
			fmt.Printf("Code %s\n", code)
		case "HOTP":
			c.log().Error(nil, "comming soon...")
		default:
			c.log().Error(nil, "secret can not be displayed")
		}
	}
}

func (c *CLI) delete(name string) {
	c.storage.DeleteSecret(name)
}

func (c *CLI) status() {
	fmt.Println("todo")
}

func (c *CLI) view() {
	cfg, err := json.MarshalIndent(c.config, "", "  ")
	if err != nil {
		c.log().Error(err, "config marshal failed")
		return
	}
	fmt.Println(string(cfg))
}

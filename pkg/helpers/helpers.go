package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/t1mon-ggg/gophkeeper/pkg/logging/zerolog"
	"golang.org/x/term"
)

var (
	termState *term.State
	cmds      []string = []string{
		"get",
		"roster",
		"revoke",
		"confirm",
		"list",
		"insert",
		"delete",
		"view",
		"edit",
		"status",
		"rollback",
		"timemachine",
	}
)

// GenSecretKey - generates a random cryptographic sequence of bytes
//		n - size of slice []byte{}
func GenSecretKey(n int) ([]byte, error) {
	data := make([]byte, n)
	_, err := rand.Read(data)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

// FileExists - check file exist or not
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GenHash - generate hashsum of []byte
func GenHash(content []byte) string {
	h := sha256.New()
	h.Write(content)
	hash := h.Sum(nil)
	return fmt.Sprintf("%x", hash)
}

// CompareHash - compare hash string with []byte
func CompareHash(hashed string, content []byte) bool {
	h := sha256.New()
	h.Write(content)
	hash := h.Sum(nil)
	return strings.EqualFold(hashed, fmt.Sprintf("%x", hash))
}

// ReadSecret - read secret from stdin in security mode
func ReadSecret(msg string) (string, error) {
	fmt.Print(msg)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	secret := string(bytePassword)
	return secret, nil
}

func SaveTermState() {
	oldState, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		return
	}
	termState = oldState
}

func RestoreTermState() {
	if termState != nil {
		term.Restore(int(os.Stdin.Fd()), termState)
	}
}

func FindCommand(in string) (string, bool) {
	for _, cmd := range cmds {
		if strings.Contains(in, cmd) {
			return cmd, true
		}
	}
	return "", false
}

func GetNameFromToken(token string) (string, error) {
	log := zerolog.New().WithPrefix("helper")
	type p struct {
		Name string `json:"name"`
		Exp  int64  `json:"exp"`
	}
	tt := strings.Split(token, ".")
	if len(tt) != 3 {
		log.Debug(nil, "token parse error")
		return "", errors.New("invalid token")
	}
	payload, err := base64.RawStdEncoding.DecodeString(tt[1])
	if err != nil {
		log.Debug(err, "base64 decode error")
		return "", err
	}
	u := new(p)
	err = json.Unmarshal(payload, u)
	if err != nil {
		log.Debug(err, "json unmarshal error")
		return "", err
	}
	return u.Name, nil
}

func GetExpirationFromToken(token string) (*time.Time, error) {
	log := zerolog.New().WithPrefix("helper")
	type p struct {
		Name string `json:"name"`
		Exp  int64  `json:"exp"`
	}
	tt := strings.Split(token, ".")
	if len(tt) != 3 {
		log.Debug(nil, "token parse error")
		return nil, errors.New("invalid token")
	}
	payload, err := base64.RawStdEncoding.DecodeString(tt[1])
	if err != nil {
		log.Debug(err, "base64 decode error")
		return nil, err
	}
	u := new(p)
	err = json.Unmarshal(payload, u)
	if err != nil {
		log.Debug(err, "json unmarshal error")
		return nil, err
	}
	t := time.Unix(u.Exp, 0)
	return &t, nil
}

// IsFlagPassed - checking the using of the flag
func IsFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

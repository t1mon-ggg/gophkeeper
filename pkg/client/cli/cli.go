package cli

import (
	"os"
	"sync"
	"time"

	prompt "github.com/c-bata/go-prompt"

	"github.com/t1mon-ggg/gophkeeper/pkg/client/config"
	"github.com/t1mon-ggg/gophkeeper/pkg/client/openpgp"
	"github.com/t1mon-ggg/gophkeeper/pkg/client/remote"
	"github.com/t1mon-ggg/gophkeeper/pkg/client/remote/web"
	"github.com/t1mon-ggg/gophkeeper/pkg/client/storage"
	"github.com/t1mon-ggg/gophkeeper/pkg/helpers"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging/zerolog"
)

type CLI struct {
	wg      *sync.WaitGroup
	storage storage.Storage
	config  *config.Config
	crypto  *openpgp.OpenPGP
	logger  logging.Logger
	api     remote.Actions
}

var livePrefixState struct {
	livePrefix string
	isEnable   bool
}

var (
	suggests       map[string][]prompt.Suggest
	s              []prompt.Suggest
	activeModeUser = []prompt.Suggest{
		{Text: "roster", Description: "list all any time logged in users"},
		{Text: "revoke", Description: "revoke pgp public key"},
		{Text: "confirm", Description: "confirm user connection and add pgp public key to key ring"},
		{Text: "quit", Description: "save changes and exit"},
		{Text: "save", Description: "save changes"},
		{Text: "..", Description: "go to up level"},
	}
	activeMode = []prompt.Suggest{
		{Text: "cmd", Description: "working area"},
		{Text: "config", Description: "configuration area"},
		{Text: "status", Description: "get current connection state"},
		{Text: "history", Description: "vault versions in time"},
		{Text: "user", Description: "setup user connections"},
		{Text: "quit", Description: "save changes and exit"},
		{Text: "save", Description: "save changes"},
	}
	activeModeHistory = []prompt.Suggest{
		{Text: "timemachine", Description: "print all time stamps of vault"},
		{Text: "rollback", Description: "rollback to vault hash"},
		{Text: "quit", Description: "save changes and exit"},
		{Text: "save", Description: "save changes"},
		{Text: "..", Description: "go to up level"},
	}
)

func init() {
	suggests = initSuggest
	s = suggests[">>> "]
}

func (c *CLI) Start() {
	if c.config.Mode != "standalone" {
		err := c.Remote()
		if err != nil {
			c.log().Error(err, "remote connection failed")
			c.log().Warn(err, "continue in standalone mode")
			c.config.Mode = "standalone"
		}
	}
	if c.config.Mode != "standalone" {
		suggests["user> "] = activeModeUser
		suggests[""] = activeMode
		suggests[">>> "] = activeMode
		suggests["history> "] = activeModeHistory
	}
	p := prompt.New(
		c.executor,
		c.completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionLivePrefix(changelivePrefix),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
	)
	p.Run()
}

func (c *CLI) log() logging.Logger {
	return c.logger
}

func New(wg *sync.WaitGroup) *CLI {
	cli := new(CLI)
	cli.logger = zerolog.New().WithPrefix("cli")
	cli.wg = wg
	cli.config = config.New()
	cli.storage = storage.New()
	cli.crypto, _ = openpgp.New()
	if cli.config.Mode != "standalone" {
		cli.api = web.New()
	}
	return cli
}

func (c *CLI) Remote() error {
	var signup bool
	err := c.api.Login(c.config.Username, c.config.Password, c.crypto.GetPublicKey())
	if err != nil {
		if err.Error() == "internal server error" || err.Error() == "bad request" {
			helpers.RestoreTermState()
			os.Exit(1)
		}
		c.log().Warn(err, "authorization failed. Try to Signup")
		err := c.api.Register(c.config.Username, c.config.Password, c.crypto.GetPublicKey())
		if err != nil {
			c.log().Error(err, "registration failed. Please contact administrator")
			helpers.RestoreTermState()
			os.Exit(1)
		}
		signup = true
	}
	pgps, err := c.api.ListPGP()
	if err != nil {
		c.log().Fatal(err, "get public keys failed")
	}
	for _, key := range pgps {
		if key.Publickey == c.crypto.GetPublicKey() && !key.Confirmed {
			c.log().Error(nil, "current pgp key pair not confirmed. please wait for confirmation")
			helpers.RestoreTermState()
			os.Exit(0)
		}
		if key.Publickey == c.crypto.GetPublicKey() {
			continue
		}
		if key.Confirmed {
			c.crypto.AddPublicKey([]byte(key.Publickey))
		}
	}
	go func() {
		err := c.api.NewStream()
		if err != nil {
			c.log().Fatal(err, "websocket connection failed")
		}
	}()
	if !signup {
		list, err := c.api.Versions()
		if err != nil {
			c.log().Fatal(err, "get list of version failed")
			return err
		}
		if len(list) == 0 {
			c.log().Info(nil, "current version is the latest")
			return nil
		}
		var latest string
		var lt time.Time
		for _, version := range list {
			if version.Date.After(lt) {
				lt = version.Date
				latest = version.Hash
			}
		}
		c.log().Trace(nil, "current version is ", c.storage.HashSum())
		c.log().Trace(nil, "latest version hash is ", latest)
		if latest != c.storage.HashSum() {
			c.log().Debug(nil, "hash sum missmatch. reloading to newest")
			body, err := c.api.Pull(latest)
			if err != nil {
				c.log().Fatal(err, "pull latest failed")
			}
			secret, err := c.crypto.DecryptWithKey(body)
			if err != nil {
				c.log().Fatal(err, "decrypt new body failed")
			}
			err = c.storage.Load(secret)
			if err != nil {
				c.log().Fatal(err, "load new version failed")
			}
		}
	}
	return nil
}

package cli

import (
	"sync"

	prompt "github.com/c-bata/go-prompt"

	"github.com/t1mon-ggg/gophkeeper/pkg/client/config"
	"github.com/t1mon-ggg/gophkeeper/pkg/client/openpgp"
	"github.com/t1mon-ggg/gophkeeper/pkg/client/storage"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging/zerolog"
)

type CLI struct {
	wg      *sync.WaitGroup
	storage storage.Storage
	config  *config.Config
	crypto  *openpgp.OpenPGP
	logger  logging.Logger
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
		{Text: "user", Description: "setup user connections"},
		{Text: "version", Description: "get binary version info"},
		{Text: "quit", Description: "save changes and exit"},
		{Text: "save", Description: "save changes"},
	}
)

func init() {
	suggests = initSuggest
	s = suggests[">>> "]
}

func (c *CLI) Start() {
	if true {
		suggests["user> "] = activeModeUser
		suggests[""] = activeMode
		suggests[">>> "] = activeMode
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
	return cli
}

package cli

import (
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/t1mon-ggg/gophkeeper/pkg/client/config"
	mockOpenPGP "github.com/t1mon-ggg/gophkeeper/pkg/client/openpgp/mock_openpgp"
	mockActions "github.com/t1mon-ggg/gophkeeper/pkg/client/remote/mock_actions"
	mockStorage "github.com/t1mon-ggg/gophkeeper/pkg/client/storage/mock_storage"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging/zerolog"
	"github.com/t1mon-ggg/gophkeeper/pkg/models"
)

func TestUserInput(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	db := mockStorage.NewMockStorage(ctl)
	actions := mockActions.NewMockActions(ctl)
	crypto := mockOpenPGP.NewMockOPENPGP(ctl)

	var errDummy = errors.New("dummy")

	gomock.InOrder(

		// actions.EXPECT().Pull("123").Return(nil, errDummy),

		// actions.EXPECT().Pull("123").Return([]byte("hello"), nil),
		// crypto.EXPECT().DecryptWithKeys([]byte("hello")).Return(nil, errDummy),

		// actions.EXPECT().Pull("123").Return([]byte("hello"), nil),
		// crypto.EXPECT().DecryptWithKeys([]byte("hello")).Return([]byte("hello"), nil),
		// db.EXPECT().Load([]byte("hello")).Return(errDummy),

		// actions.EXPECT().Pull("123").Return([]byte("hello"), nil),
		// crypto.EXPECT().DecryptWithKeys([]byte("hello")).Return([]byte("hello"), nil),
		// db.EXPECT().Load([]byte("hello")).Return(nil),

		// save actions
		db.EXPECT().Save().Return([]byte("hello"), nil),
		crypto.EXPECT().EncryptWithKeys([]byte("hello")).Return([]byte("hello"), nil),
		db.EXPECT().HashSum().Return("1234"),
		actions.EXPECT().Push("hello", "1234").Return(nil),

		db.EXPECT().Save().Return([]byte("hello"), nil),
		crypto.EXPECT().EncryptWithKeys([]byte("hello")).Return([]byte("hello"), nil),
		db.EXPECT().HashSum().Return("1234"),
		actions.EXPECT().Push("hello", "1234").Return(nil),

		//quit action
		db.EXPECT().Save().Return([]byte("hello"), nil),
		crypto.EXPECT().EncryptWithKeys([]byte("hello")).Return([]byte("hello"), nil),
		db.EXPECT().HashSum().Return("1234"),
		actions.EXPECT().Push("hello", "1234").Return(nil),
		actions.EXPECT().Close().Return(nil),

		// timemachine action
		actions.EXPECT().Versions().Return([]models.Version{}, nil),

		actions.EXPECT().Versions().Return(nil, errDummy),

		actions.EXPECT().Versions().Return([]models.Version{{Date: time.Now(), Hash: "123"}}, nil),

		// roster action
		actions.EXPECT().ListPGP().Return(nil, errDummy),

		actions.EXPECT().ListPGP().Return([]models.PGP{}, nil),

		//revoke action
		crypto.EXPECT().GetPublicKey().Return("321"),

		crypto.EXPECT().GetPublicKey().Return("321"),
		actions.EXPECT().ListPGP().Return(nil, errDummy),

		crypto.EXPECT().GetPublicKey().Return("321"),
		actions.EXPECT().ListPGP().Return([]models.PGP{{Date: time.Now(), Publickey: "KEY", Confirmed: true}, {Date: time.Now(), Publickey: "KEY2", Confirmed: true}}, nil),
		actions.EXPECT().RevokePGP("KEY").Return(errDummy),

		crypto.EXPECT().GetPublicKey().Return("321"),
		actions.EXPECT().ListPGP().Return([]models.PGP{{Date: time.Now(), Publickey: "KEY", Confirmed: true}, {Date: time.Now().Add(time.Hour), Publickey: "KEY2", Confirmed: true}}, nil),
		actions.EXPECT().RevokePGP("KEY").Return(nil),
		crypto.EXPECT().ReloadPublicKeys([]string{"KEY2"}).Return(errDummy),

		crypto.EXPECT().GetPublicKey().Return("321"),
		actions.EXPECT().ListPGP().Return([]models.PGP{{Date: time.Now(), Publickey: "KEY", Confirmed: true}, {Date: time.Now().Add(time.Hour), Publickey: "KEY2", Confirmed: true}}, nil),
		actions.EXPECT().RevokePGP("KEY").Return(nil),
		crypto.EXPECT().ReloadPublicKeys([]string{"KEY2"}).Return(nil),
		db.EXPECT().ReEncrypt().Return(nil, errDummy),

		crypto.EXPECT().GetPublicKey().Return("321"),
		actions.EXPECT().ListPGP().Return([]models.PGP{{Date: time.Now(), Publickey: "KEY", Confirmed: true}, {Date: time.Now().Add(time.Hour), Publickey: "KEY2", Confirmed: true}}, nil),
		actions.EXPECT().RevokePGP("KEY").Return(nil),
		crypto.EXPECT().ReloadPublicKeys([]string{"KEY2"}).Return(nil),
		db.EXPECT().ReEncrypt().Return([]byte("hello"), nil),
		crypto.EXPECT().EncryptWithKeys([]byte("hello")).Return(nil, errDummy),

		crypto.EXPECT().GetPublicKey().Return("321"),
		actions.EXPECT().ListPGP().Return([]models.PGP{{Date: time.Now(), Publickey: "KEY", Confirmed: true}, {Date: time.Now().Add(time.Hour), Publickey: "KEY2", Confirmed: true}}, nil),
		actions.EXPECT().RevokePGP("KEY").Return(nil),
		crypto.EXPECT().ReloadPublicKeys([]string{"KEY2"}).Return(nil),
		db.EXPECT().ReEncrypt().Return([]byte("hello"), nil),
		crypto.EXPECT().EncryptWithKeys([]byte("hello")).Return([]byte("hello"), nil),
		db.EXPECT().HashSum().Return("123"),
		actions.EXPECT().Push("hello", "123").Return(errDummy),

		crypto.EXPECT().GetPublicKey().Return("321"),
		actions.EXPECT().ListPGP().Return([]models.PGP{{Date: time.Now(), Publickey: "KEY", Confirmed: true}, {Date: time.Now().Add(time.Hour), Publickey: "KEY2", Confirmed: true}}, nil),
		actions.EXPECT().RevokePGP("KEY").Return(nil),
		crypto.EXPECT().ReloadPublicKeys([]string{"KEY2"}).Return(nil),
		db.EXPECT().ReEncrypt().Return([]byte("hello"), nil),
		crypto.EXPECT().EncryptWithKeys([]byte("hello")).Return([]byte("hello"), nil),
		db.EXPECT().HashSum().Return("123"),
		actions.EXPECT().Push("hello", "123").Return(nil),

		//confirm action
		actions.EXPECT().ListPGP().Return(nil, errDummy),

		actions.EXPECT().ListPGP().Return([]models.PGP{}, nil),

		actions.EXPECT().ListPGP().Return([]models.PGP{{Date: time.Now(), Publickey: "KEY", Confirmed: false}, {Date: time.Now().Add(time.Hour), Publickey: "KEY2", Confirmed: true}}, nil),
		actions.EXPECT().ConfirmPGP("KEY").Return(errDummy),

		actions.EXPECT().ListPGP().Return([]models.PGP{{Date: time.Now(), Publickey: "KEY", Confirmed: false}, {Date: time.Now().Add(time.Hour), Publickey: "KEY2", Confirmed: true}}, nil),
		actions.EXPECT().ConfirmPGP("KEY").Return(nil),
		crypto.EXPECT().AddPublicKey([]byte("KEY")).Return(errDummy),

		actions.EXPECT().ListPGP().Return([]models.PGP{{Date: time.Now(), Publickey: "KEY", Confirmed: false}, {Date: time.Now().Add(time.Hour), Publickey: "KEY2", Confirmed: true}}, nil),
		actions.EXPECT().ConfirmPGP("KEY").Return(nil),
		crypto.EXPECT().AddPublicKey([]byte("KEY")).Return(nil),
		db.EXPECT().ReEncrypt().Return(nil, errDummy),

		actions.EXPECT().ListPGP().Return([]models.PGP{{Date: time.Now(), Publickey: "KEY", Confirmed: false}, {Date: time.Now().Add(time.Hour), Publickey: "KEY2", Confirmed: true}}, nil),
		actions.EXPECT().ConfirmPGP("KEY").Return(nil),
		crypto.EXPECT().AddPublicKey([]byte("KEY")).Return(nil),
		db.EXPECT().ReEncrypt().Return([]byte("hello"), nil),
		crypto.EXPECT().EncryptWithKeys([]byte("hello")).Return(nil, errDummy),

		actions.EXPECT().ListPGP().Return([]models.PGP{{Date: time.Now(), Publickey: "KEY", Confirmed: false}, {Date: time.Now().Add(time.Hour), Publickey: "KEY2", Confirmed: true}}, nil),
		actions.EXPECT().ConfirmPGP("KEY").Return(nil),
		crypto.EXPECT().AddPublicKey([]byte("KEY")).Return(nil),
		db.EXPECT().ReEncrypt().Return([]byte("hello"), nil),
		crypto.EXPECT().EncryptWithKeys([]byte("hello")).Return([]byte("hello"), nil),
		db.EXPECT().HashSum().Return("123"),
		actions.EXPECT().Push("hello", "123").Return(errDummy),

		actions.EXPECT().ListPGP().Return([]models.PGP{{Date: time.Now(), Publickey: "KEY", Confirmed: false}, {Date: time.Now().Add(time.Hour), Publickey: "KEY2", Confirmed: true}}, nil),
		actions.EXPECT().ConfirmPGP("KEY").Return(nil),
		crypto.EXPECT().AddPublicKey([]byte("KEY")).Return(nil),
		db.EXPECT().ReEncrypt().Return([]byte("hello"), nil),
		crypto.EXPECT().EncryptWithKeys([]byte("hello")).Return([]byte("hello"), nil),
		db.EXPECT().HashSum().Return("123"),
		actions.EXPECT().Push("hello", "123").Return(nil),

		// rollback action

		// actions.EXPECT().ListPGP().Return([]models.PGP{{Date: time.Now(), Publickey: "KEY", Confirmed: false}, {Date: time.Now().Add(time.Hour), Publickey: "KEY2", Confirmed: true}}, nil),
		// actions.EXPECT().ConfirmPGP("KEY").Return(nil),
	)

	cli := &CLI{
		wg:      new(sync.WaitGroup),
		storage: db,
		logger:  zerolog.New().WithPrefix("operations-test"),
		api:     actions,
		config:  config.New(),
		crypto:  crypto,
	}
	cli.config.Mode = "client-server"
	cli.wg.Add(1)

	type test struct {
		name       string
		in         string
		livePrefix string
	}

	tests := []test{
		{
			name: "empty action",
			in:   "",
		},
		{
			name:       "save action",
			in:         "save",
			livePrefix: ">>> ",
		},
		{
			name:       "save action",
			in:         "save",
			livePrefix: "",
		},
		{
			name:       "quit action",
			in:         "quit",
			livePrefix: "",
		},
		{
			name:       "up action",
			in:         "..",
			livePrefix: ">>> ",
		},
		{
			name:       "up action",
			in:         "..",
			livePrefix: "cmd> ",
		},
		{
			name:       "up action",
			in:         "..",
			livePrefix: "cmd/test> ",
		},
		{
			name:       "timemachine action",
			in:         "timemachine",
			livePrefix: "history> ",
		},
		{
			name:       "timemachine action",
			in:         "timemachine",
			livePrefix: "history> ",
		},
		{
			name:       "timemachine action",
			in:         "timemachine",
			livePrefix: "history> ",
		},
		{
			name:       "timemachine action",
			in:         "timemachine",
			livePrefix: "cmd> ",
		},
		{
			name:       "roster action",
			in:         "roster",
			livePrefix: "cmd> ",
		},
		{
			name:       "roster action",
			in:         "roster",
			livePrefix: "user> ",
		},
		{
			name:       "roster action",
			in:         "roster",
			livePrefix: "user> ",
		},
		{
			name:       "roster action",
			in:         "roster",
			livePrefix: "history> ",
		},
		{
			name:       "revoke action",
			in:         "revoke",
			livePrefix: "user> ",
		},
		{
			name:       "revoke action",
			in:         "revoke 8d23cf6c86e834a7aa6eded54c26ce2bb2e74903538c61bdd5d2197997ab2f72",
			livePrefix: "user> ",
		},
		{
			name:       "revoke action",
			in:         "revoke a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3",
			livePrefix: "user> ",
		},
		{
			name:       "revoke action",
			in:         "revoke 5ca24005b740717ba4f3f6bc48a230700e68c2a4b11ecedb96f169f4efaf1f21",
			livePrefix: "user> ",
		},
		{
			name:       "revoke action",
			in:         "revoke 5ca24005b740717ba4f3f6bc48a230700e68c2a4b11ecedb96f169f4efaf1f21",
			livePrefix: "user> ",
		},
		{
			name:       "revoke action",
			in:         "revoke 5ca24005b740717ba4f3f6bc48a230700e68c2a4b11ecedb96f169f4efaf1f21",
			livePrefix: "user> ",
		},
		{
			name:       "revoke action",
			in:         "revoke 5ca24005b740717ba4f3f6bc48a230700e68c2a4b11ecedb96f169f4efaf1f21",
			livePrefix: "user> ",
		},
		{
			name:       "revoke action",
			in:         "revoke 5ca24005b740717ba4f3f6bc48a230700e68c2a4b11ecedb96f169f4efaf1f21",
			livePrefix: "user> ",
		},
		{
			name:       "revoke action",
			in:         "revoke 5ca24005b740717ba4f3f6bc48a230700e68c2a4b11ecedb96f169f4efaf1f21",
			livePrefix: "user> ",
		},
		{
			name:       "revoke action",
			in:         "revoke 5ca24005b740717ba4f3f6bc48a230700e68c2a4b11ecedb96f169f4efaf1f21",
			livePrefix: "cmd> ",
		},
		{
			name:       "confirm action",
			in:         "confirm 5ca24005b740717ba4f3f6bc48a230700e68c2a4b11ecedb96f169f4efaf1f21",
			livePrefix: "cmd> ",
		},
		{
			name:       "confirm action",
			in:         "confirm 5ca24005b740717ba4f3f6bc48a230700e68c2a4b11ecedb96f169f4efaf1f21",
			livePrefix: "user> ",
		},
		{
			name:       "confirm action",
			in:         "confirm 5ca24005b740717ba4f3f6bc48a230700e68c2a4b11ecedb96f169f4efaf1f21",
			livePrefix: "user> ",
		},
		{
			name:       "confirm action",
			in:         "confirm 5ca24005b740717ba4f3f6bc48a230700e68c2a4b11ecedb96f169f4efaf1f21",
			livePrefix: "user> ",
		},
		{
			name:       "confirm action",
			in:         "confirm 5ca24005b740717ba4f3f6bc48a230700e68c2a4b11ecedb96f169f4efaf1f21",
			livePrefix: "user> ",
		},
		{
			name:       "confirm action",
			in:         "confirm 5ca24005b740717ba4f3f6bc48a230700e68c2a4b11ecedb96f169f4efaf1f21",
			livePrefix: "user> ",
		},
		{
			name:       "confirm action",
			in:         "confirm 5ca24005b740717ba4f3f6bc48a230700e68c2a4b11ecedb96f169f4efaf1f21",
			livePrefix: "user> ",
		},
		{
			name:       "confirm action",
			in:         "confirm 5ca24005b740717ba4f3f6bc48a230700e68c2a4b11ecedb96f169f4efaf1f21",
			livePrefix: "user> ",
		},
		{
			name:       "confirm action",
			in:         "confirm 5ca24005b740717ba4f3f6bc48a230700e68c2a4b11ecedb96f169f4efaf1f21",
			livePrefix: "user> ",
		},
	}
	livePrefixState.isEnable = true
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			livePrefixState.livePrefix = tt.livePrefix
			cli.executor(tt.in)
			os.Remove("secrets.db")
		})
	}

}

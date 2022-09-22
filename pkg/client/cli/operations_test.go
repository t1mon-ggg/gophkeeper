package cli

import (
	"errors"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/t1mon-ggg/gophkeeper/pkg/client/config"
	mockOpenPGP "github.com/t1mon-ggg/gophkeeper/pkg/client/openpgp/mock_openpgp"
	mockActions "github.com/t1mon-ggg/gophkeeper/pkg/client/remote/mock_actions"
	"github.com/t1mon-ggg/gophkeeper/pkg/client/storage"
	mockStorage "github.com/t1mon-ggg/gophkeeper/pkg/client/storage/mock_storage"
	"github.com/t1mon-ggg/gophkeeper/pkg/helpers"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging/zerolog"
)

func TestSave(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	db := mockStorage.NewMockStorage(ctl)
	actions := mockActions.NewMockActions(ctl)
	crypto := mockOpenPGP.NewMockOPENPGP(ctl)

	gomock.InOrder(
		db.EXPECT().Save().Return([]byte("hello"), nil),
		crypto.EXPECT().EncryptWithKeys([]byte("hello")).Return([]byte("hello"), nil),
		db.EXPECT().HashSum().Return("1234"),
		actions.EXPECT().Push("hello", "1234").Return(nil),

		db.EXPECT().Save().Return(nil, storage.ErrHashValid),

		db.EXPECT().Save().Return(nil, errors.New("dummy")),

		db.EXPECT().Save().Return([]byte("hello"), nil),

		db.EXPECT().Save().Return([]byte("hello"), nil),
		crypto.EXPECT().EncryptWithKeys([]byte("hello")).Return(nil, errors.New("dummy")),

		db.EXPECT().Save().Return([]byte("hello"), nil),
		crypto.EXPECT().EncryptWithKeys([]byte("hello")).Return([]byte("hello"), nil),
		db.EXPECT().HashSum().Return("1234"),
		actions.EXPECT().Push("hello", "1234").Return(errors.New("dummy")),
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
	cli.save()
	f, err := os.Open("secrets.db")
	require.NoError(t, err)
	buf, err := io.ReadAll(f)
	require.NoError(t, err)

	require.Equal(t, []byte("hello"), buf)
	err = os.Remove("secrets.db")
	require.NoError(t, err)

	cli.save()
	require.False(t, helpers.FileExists("secrets.db"))

	cli.save()
	require.False(t, helpers.FileExists("secrets.db"))

	f, err = os.OpenFile("secrets.db", os.O_CREATE, 0000)
	require.NoError(t, err)
	cli.save()
	f.Close()
	_, err = os.Open("secrets.db")
	require.Error(t, err)
	err = os.Remove("secrets.db")
	require.NoError(t, err)

	cli.save()
	f, err = os.Open("secrets.db")
	require.NoError(t, err)
	buf, err = io.ReadAll(f)
	require.NoError(t, err)
	require.Equal(t, []byte{}, buf)
	err = os.Remove("secrets.db")
	require.NoError(t, err)

	cli.save()
}

func TestInsert(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	db := mockStorage.NewMockStorage(ctl)
	actions := mockActions.NewMockActions(ctl)
	crypto := mockOpenPGP.NewMockOPENPGP(ctl)

	gomock.InOrder()

	cli := &CLI{
		wg:      new(sync.WaitGroup),
		storage: db,
		logger:  zerolog.New().WithPrefix("operations-test"),
		api:     actions,
		config:  config.New(),
		crypto:  crypto,
	}
	cli.config.Mode = "client-server"

}

func TestList(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	db := mockStorage.NewMockStorage(ctl)
	actions := mockActions.NewMockActions(ctl)
	crypto := mockOpenPGP.NewMockOPENPGP(ctl)

	gomock.InOrder()

	cli := &CLI{
		wg:      new(sync.WaitGroup),
		storage: db,
		logger:  zerolog.New().WithPrefix("operations-test"),
		api:     actions,
		config:  config.New(),
		crypto:  crypto,
	}
	cli.config.Mode = "client-server"
}

func TestGet(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	db := mockStorage.NewMockStorage(ctl)
	actions := mockActions.NewMockActions(ctl)
	crypto := mockOpenPGP.NewMockOPENPGP(ctl)

	gomock.InOrder()

	cli := &CLI{
		wg:      new(sync.WaitGroup),
		storage: db,
		logger:  zerolog.New().WithPrefix("operations-test"),
		api:     actions,
		config:  config.New(),
		crypto:  crypto,
	}
	cli.config.Mode = "client-server"
}

func TestDelete(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	db := mockStorage.NewMockStorage(ctl)
	actions := mockActions.NewMockActions(ctl)
	crypto := mockOpenPGP.NewMockOPENPGP(ctl)

	gomock.InOrder()

	cli := &CLI{
		wg:      new(sync.WaitGroup),
		storage: db,
		logger:  zerolog.New().WithPrefix("operations-test"),
		api:     actions,
		config:  config.New(),
		crypto:  crypto,
	}
	cli.config.Mode = "client-server"
}

func TestStatus(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	db := mockStorage.NewMockStorage(ctl)
	actions := mockActions.NewMockActions(ctl)
	crypto := mockOpenPGP.NewMockOPENPGP(ctl)

	gomock.InOrder()

	cli := &CLI{
		wg:      new(sync.WaitGroup),
		storage: db,
		logger:  zerolog.New().WithPrefix("operations-test"),
		api:     actions,
		config:  config.New(),
		crypto:  crypto,
	}
	cli.config.Mode = "client-server"
}

func TestView(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	db := mockStorage.NewMockStorage(ctl)
	actions := mockActions.NewMockActions(ctl)
	crypto := mockOpenPGP.NewMockOPENPGP(ctl)

	gomock.InOrder()

	cli := &CLI{
		wg:      new(sync.WaitGroup),
		storage: db,
		logger:  zerolog.New().WithPrefix("operations-test"),
		api:     actions,
		config:  config.New(),
		crypto:  crypto,
	}
	cli.config.Mode = "client-server"
}

func TestConfirm(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	db := mockStorage.NewMockStorage(ctl)
	actions := mockActions.NewMockActions(ctl)
	crypto := mockOpenPGP.NewMockOPENPGP(ctl)

	gomock.InOrder()

	cli := &CLI{
		wg:      new(sync.WaitGroup),
		storage: db,
		logger:  zerolog.New().WithPrefix("operations-test"),
		api:     actions,
		config:  config.New(),
		crypto:  crypto,
	}
	cli.config.Mode = "client-server"
}
func TestRevoke(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	db := mockStorage.NewMockStorage(ctl)
	actions := mockActions.NewMockActions(ctl)
	crypto := mockOpenPGP.NewMockOPENPGP(ctl)

	gomock.InOrder()

	cli := &CLI{
		wg:      new(sync.WaitGroup),
		storage: db,
		logger:  zerolog.New().WithPrefix("operations-test"),
		api:     actions,
		config:  config.New(),
		crypto:  crypto,
	}
	cli.config.Mode = "client-server"
}

func TestRoster(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	db := mockStorage.NewMockStorage(ctl)
	actions := mockActions.NewMockActions(ctl)
	crypto := mockOpenPGP.NewMockOPENPGP(ctl)

	gomock.InOrder()

	cli := &CLI{
		wg:      new(sync.WaitGroup),
		storage: db,
		logger:  zerolog.New().WithPrefix("operations-test"),
		api:     actions,
		config:  config.New(),
		crypto:  crypto,
	}
	cli.config.Mode = "client-server"
}

func TestRollback(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	db := mockStorage.NewMockStorage(ctl)
	actions := mockActions.NewMockActions(ctl)
	crypto := mockOpenPGP.NewMockOPENPGP(ctl)

	gomock.InOrder()

	cli := &CLI{
		wg:      new(sync.WaitGroup),
		storage: db,
		logger:  zerolog.New().WithPrefix("operations-test"),
		api:     actions,
		config:  config.New(),
		crypto:  crypto,
	}
	cli.config.Mode = "client-server"
}

func TestTimemachine(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	db := mockStorage.NewMockStorage(ctl)
	actions := mockActions.NewMockActions(ctl)
	crypto := mockOpenPGP.NewMockOPENPGP(ctl)

	gomock.InOrder()

	cli := &CLI{
		wg:      new(sync.WaitGroup),
		storage: db,
		logger:  zerolog.New().WithPrefix("operations-test"),
		api:     actions,
		config:  config.New(),
		crypto:  crypto,
	}
	cli.config.Mode = "client-server"
}

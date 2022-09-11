package storage

import (
	"net"
	"time"

	"github.com/t1mon-ggg/gophkeeper/pkg/logging"
	"github.com/t1mon-ggg/gophkeeper/pkg/models"
)

// Storage
type Storage interface {
	SignUp(username, password string) error
	SignIn(user models.User) error
	DeleteUser(username string) error

	Push(username, checksum string, data []byte) error
	Pull(username, checksum string) ([]byte, error)
	SaveLog(username, action, checksum, sign string, ip *net.IP, date time.Time) error
	GetLog(username string) ([]models.Action, error)
	AddPGP(username, publickey string) error
	DeletePGP(publickey string) error

	Close()
	Ping() error

	log() logging.Logger
}

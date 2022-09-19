package storage

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io"
	"strings"
	"sync"

	"github.com/mgutz/ansi"

	"github.com/t1mon-ggg/gophkeeper/pkg/client/storage/secrets"
	"github.com/t1mon-ggg/gophkeeper/pkg/helpers"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging/zerolog"
)

var (
	ErrHashValid = errors.New("hashsum is identical")
	once         sync.Once
	_storage     Storage
)

func init() {
	gob.Register(&secrets.AnyBinary{})
	gob.Register(&secrets.AnyText{})
	gob.Register(&secrets.CreditCard{})
	gob.Register(&secrets.OTP{})
	gob.Register(&secrets.UserPass{})
}

type Storage interface {
	Save() ([]byte, error)
	Load(b []byte) error
	InsertSecret(name, description string, secret Secret) Storage
	DeleteSecret(name string) Storage
	GetSecret(name string) Secret
	ListSecrets() map[string]string
	HashSum() string
}
type Secret interface {
	Scope() string
	Value() any
}

type keeper struct {
	secrets []Value
	logger  logging.Logger
	rwMutex *sync.RWMutex
	hashsum string
}

type Value struct {
	Name        string
	Description string
	Record      Secret
}

func New() Storage {
	once.Do(func() {
		k := new(keeper)
		k.logger = zerolog.New().WithPrefix("storage")
		k.rwMutex = new(sync.RWMutex)
		k.secrets = []Value{}
		_storage = k
	})
	return _storage
}

func (k *keeper) HashSum() string {
	return k.hashsum
}

func (k *keeper) Save() ([]byte, error) {
	if len(k.secrets) == 0 {
		return []byte{}, nil
	}
	k.rwMutex.RLock()
	defer k.rwMutex.RUnlock()
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(k.secrets)
	if err != nil {
		return []byte{}, err
	}
	hash := helpers.GenHash(buf.Bytes())
	if hash == k.HashSum() {
		k.logger.Info(nil, "hashsum identical. skip loading")
		return []byte{}, ErrHashValid
	}
	k.hashsum = hash
	return buf.Bytes(), nil
}

func (k *keeper) Load(b []byte) error {
	hash := helpers.GenHash(b)
	if hash == k.HashSum() {
		k.logger.Info(nil, "hashsum identical. skip loading")
		return nil
	}
	k.rwMutex.Lock()
	defer k.rwMutex.Unlock()
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	v := []Value{}
	err := dec.Decode(&v)
	if err != nil {
		if err == io.EOF {
			k.secrets = []Value{}
			return nil
		}
		return err
	}
	k.hashsum = hash
	k.secrets = v
	return nil
}

func (k *keeper) InsertSecret(name, description string, secret Secret) Storage {
	list := k.ListSecrets()
	if len(list) != 0 {
		if _, ok := list[name]; ok {
			k.logger.Warnf("Create new secret failed. Such secret name %s already exists", nil, ansi.Color(name, "red+b"))
			return k
		}
	}
	k.rwMutex.Lock()
	defer k.rwMutex.Unlock()
	k.secrets = append(k.secrets,
		Value{
			Name:        name,
			Description: description,
			Record:      secret,
		})
	return k
}

func (k *keeper) DeleteSecret(name string) Storage {
	k.rwMutex.Lock()
	defer k.rwMutex.Unlock()
	kk := make([]Value, len(k.secrets)-1)
	var found bool
	for _, v := range k.secrets {
		if strings.EqualFold(name, v.Name) {
			found = true
			continue
		}
		kk = append(kk, v)
	}
	if !found {
		k.logger.Warnf("Delete secret failed. No such secret name %s", nil, ansi.Color(name, "red+b"))
		return k
	}
	k.secrets = kk
	return k
}

func (k *keeper) GetSecret(name string) Secret {
	k.rwMutex.RLock()
	defer k.rwMutex.RUnlock()
	for _, v := range k.secrets {
		if strings.EqualFold(name, v.Name) {
			return v.Record
		}
	}
	k.logger.Warnf("Get secret failed. No such secret name %s", nil, ansi.Color(name, "red+b"))
	return nil
}

func (k *keeper) ListSecrets() map[string]string {
	k.rwMutex.RLock()
	defer k.rwMutex.RUnlock()
	list := make(map[string]string, len(k.secrets))
	for _, v := range k.secrets {
		list[v.Name] = v.Description
	}
	delete(list, "")
	return list
}

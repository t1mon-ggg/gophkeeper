package storage

import (
	"bytes"
	"encoding/gob"
	"strings"
	"sync"

	"github.com/mgutz/ansi"

	"github.com/t1mon-ggg/gophkeeper/pkg/client/storage/secrets"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging"
)

func init() {
	gob.Register(&secrets.AnyBinary{})
	gob.Register(&secrets.AnyText{})
	gob.Register(&secrets.CreditCard{})
	gob.Register(&secrets.OTP{})
	gob.Register(&secrets.UserPass{})
}

type keeper struct {
	secrets []Value
	logger  logging.Logger
	rwMutex *sync.RWMutex
}

type Value struct {
	Name        string
	Description string
	Record      Secret
}

type Secret interface {
	Scope() string
	Value() any
}

func New(logger logging.Logger) *keeper {
	k := new(keeper)
	k.logger = logger.WithPrefix("storage")
	k.rwMutex = new(sync.RWMutex)
	k.secrets = []Value{}
	return k
}

func (k *keeper) Save() ([]byte, error) {
	k.rwMutex.RLock()
	defer k.rwMutex.RUnlock()
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(k.secrets)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func (k *keeper) Load(b []byte) error {
	k.rwMutex.Lock()
	defer k.rwMutex.Unlock()
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	v := []Value{}
	err := dec.Decode(&v)
	if err != nil {
		return err
	}
	k.secrets = v
	return nil
}

func (k *keeper) InsertSecret(name, description string, secret Secret) *keeper {
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

func (k *keeper) DeleteSecret(name string) *keeper {
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

func (k *keeper) Log() logging.Logger {
	return k.logger
}

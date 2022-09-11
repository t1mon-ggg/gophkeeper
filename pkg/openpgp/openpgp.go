package openpgp

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"time"

	pgpcrypto "github.com/ProtonMail/gopenpgp/v2/crypto"

	"github.com/t1mon-ggg/gophkeeper/pkg/helpers"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging/zerolog"
)

type PGPFile struct {
	filename string
	keyring  *OpenPGP
	data     []byte
	buf      *bytes.Buffer
}

var (
	log    logging.Logger
	errPGP = errors.New("openpgp failed")
)

func init() {
	log = zerolog.New().WithPrefix("openpgp")
}

type OpenPGP struct {
	Public     *pgpcrypto.KeyRing
	Private    *pgpcrypto.KeyRing
	passphrase []byte
}

func New(passphrase string) (*OpenPGP, error) {
	var err error
	p := new(OpenPGP)
	p.passphrase = []byte(passphrase)
	p.Private, err = pgpcrypto.NewKeyRing(nil)
	if err != nil {
		log.Debug(err, "empty private keyring creation failed")
		return nil, errPGP
	}
	p.Public, err = pgpcrypto.NewKeyRing(nil)
	if err != nil {
		log.Debug(err, "empty public keyring creation failed")
		return nil, errPGP
	}
	return p, nil
}

func (p *OpenPGP) AddPublicKey(armored []byte) error {
	key, err := pgpcrypto.NewKeyFromArmored(string(armored))
	if err != nil {
		log.Debug(err, "key parse failed")
		return errPGP
	}
	err = p.Public.AddKey(key)
	if err != nil {
		log.Debug(err, "key add failed")
		return errPGP
	}
	log.Debug(nil, "key added")
	return nil
}

func (p *OpenPGP) AddPrivateKey(armored []byte) error {
	key, err := pgpcrypto.NewKeyFromArmored(string(armored))
	if err != nil {
		log.Debug(err, "key parse failed")
		return errPGP
	}
	ok, err := key.IsLocked()
	if err != nil {
		log.Debug(err, "key locked state check failed")
		return errPGP
	}
	if ok {
		unlocked, err := key.Unlock(p.passphrase)
		if err != nil {
			log.Debug(err, "key unlock failed")
			return errPGP
		}
		err = p.Private.AddKey(unlocked)
		if err != nil {
			log.Debug(err, "key add failed")
			return errPGP
		}
		log.Debug(nil, "locked key added")
		return nil
	}
	err = p.Private.AddKey(key)
	if err != nil {
		log.Debug(err, "key add failed")
		return errPGP
	}
	log.Debug(nil, "key added")
	return nil
}

func (p *OpenPGP) GeneratePair(name, email string) error {
	key, err := pgpcrypto.GenerateKey(name, email, "x25519", 0)
	if err != nil {
		log.Debug(err, "openpgp key pair generation failed")
		return errPGP
	}
	defer key.ClearPrivateParams()
	locked, err := key.Lock([]byte(p.passphrase))
	if err != nil {
		log.Debug(err, "unable to lock openpgp key")
		return errPGP
	}
	if !helpers.FileExists("./openpgp") {
		err := os.MkdirAll("./openpgp", 0700)
		if err != nil {
			log.Debug(err, "openpgp folder creation failed")
			return errPGP
		}
	}
	keyfilename := fmt.Sprintf("./openpgp/%s.gpg", name)
	pubfilename := fmt.Sprintf("./openpgp/%s.gpg.pub", name)
	keyFile, err := os.OpenFile(keyfilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Debug(err, "openpgp private key file creation failed")
		return errPGP
	}
	defer keyFile.Close()
	pubFile, err := os.OpenFile(pubfilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Debug(err, "openpgp public key file creation failed")
		return errPGP
	}
	defer pubFile.Close()
	public, err := key.GetArmoredPublicKey()
	if err != nil {
		log.Debug(err, "get openpgp public key failed")
		return errPGP
	}
	fmt.Fprint(pubFile, public)
	private, err := locked.Armor()
	if err != nil {
		log.Debug(err, "get openpgp locked private key failed")
		return errPGP
	}
	fmt.Fprint(keyFile, private)
	err = p.AddPublicKey([]byte(public))
	if err != nil {
		log.Debug(err, "add public key to keyring failed")
		return errPGP
	}
	p.AddPrivateKey([]byte(private))
	if err != nil {
		log.Debug(err, "add private key to keyring failed")
		return errPGP
	}
	return nil
}

func (p *OpenPGP) EncryptWithKeys(data []byte) (*pgpcrypto.PGPMessage, error) {
	if !p.Public.CanEncrypt() {
		log.Debug(nil, "pgp keyring can not be used for encryption")
		return nil, errPGP
	}
	msg := pgpcrypto.NewPlainMessage(data)
	encryptedMsg, err := p.Public.EncryptWithCompression(msg, p.Private)
	if err != nil {
		log.Debug(err, "data encryption failed")
		return nil, errPGP
	}
	return encryptedMsg, nil
}

func (p *OpenPGP) DecryptWithKey(data []byte) ([]byte, error) {
	encryptedMsg := pgpcrypto.NewPGPMessage(data)
	clear, err := p.Private.Decrypt(encryptedMsg, p.Public, time.Now().Unix())
	if err != nil {
		log.Debug(err, "data decryption failed")
		return nil, errPGP
	}
	return clear.Data, nil
}

// func (p *OpenPGP) NewPGPStream(filename string) (*PGPFile, error) {
// 	pgp := new(PGPFile)
// 	pgp.filename = filename
// 	pgp.keyring = p
// 	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDONLY, 0644)
// 	if err != nil {
// 		return nil, err
// 	}
// 	filebuf := bufio.NewReader(file)
// 	var buf bytes.Buffer
// 	part := make([]byte, 1024)
// 	for {
// 		count, err := filebuf.Read(part)
// 		if err != nil {
// 			if err != io.EOF {
// 				return nil, err
// 			}
// 			break
// 		}

// 		fmt.Println("!!!!!!!!!!", part)
// 		buf.Write(part[:count])
// 	}
// 	data, err := p.DecryptWithKey(buf.Bytes())
// 	if err != nil {
// 		return nil, err
// 	}
// 	pgp.data = data
// 	pgp.buf = bytes.NewBuffer(pgp.data)
// 	fmt.Println("@@@@@@@@@@@@@@@@@@", pgp.buf.Bytes())
// 	fmt.Println("@@@@@@@@@@@@@@@@@@", data)
// 	return pgp, nil
// }

// func (p *PGPFile) Close() error {
// 	encrypted, err := p.keyring.EncryptWithKeys(p.buf.Bytes())
// 	if err != nil {
// 		return err
// 	}
// 	file, err := os.OpenFile(p.filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = file.Write(encrypted.Data)
// 	if err != nil {
// 		return err
// 	}
// 	file.Close()
// 	return nil
// }

// func (p *PGPFile) Reader() io.Reader {
// 	return p.buf
// }

// func (p *PGPFile) Writer() io.Writer {
// 	return p.buf
// }

// func (p *PGPFile) ReadWriter() io.ReadWriter {
// 	x := bufio.NewReadWriter(bufio.NewReader(p.buf), bufio.NewWriter(p.buf))
// 	return x
// }

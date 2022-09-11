package openpgp

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging/zerolog"
)

func TestGeneratePair(t *testing.T) {
	type args struct {
		name       string
		email      string
		passphrase string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "genereate openpgp pair",
			args: args{
				name:       "tester",
				email:      "tester@test.com",
				passphrase: "test-password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := New(tt.args.passphrase)
			require.NoError(t, err)
			err = p.GeneratePair(tt.args.name, tt.args.email)
			require.NoError(t, err)
			log.Debugf("user %s pair created and tested", nil, tt.args.name)
		})
	}
	content, err := os.ReadDir("./openpgp")
	require.NoError(t, err)
	for _, file := range content {
		path := fmt.Sprintf("./openpgp/%s", file.Name())
		err := os.Remove(path)
		require.NoError(t, err)
		log.Debugf("%s deleted", nil, path)
	}
	err = os.RemoveAll("./openpgp")
	require.NoError(t, err)
	log.Debug(nil, "./openpgp")
}

func TestGenerateEncryptionAndDecryption(t *testing.T) {
	log := zerolog.New().WithPrefix("openpgp-test")
	type users struct {
		name       string
		email      string
		passphrase string
	}
	usrs := []users{
		{
			name:       "first",
			email:      "first@test.com",
			passphrase: "first_passphrase",
		},
		{
			name:       "second",
			email:      "second@test.com",
			passphrase: "second_passphrase",
		},
		{
			name:       "third",
			email:      "third@test.com",
			passphrase: "third_passphrase",
		},
	}
	msg := []byte("this is very confidential information")
	for _, user := range usrs {
		p, err := New(user.passphrase)
		require.NoError(t, err)
		err = p.GeneratePair(user.name, user.email)
		require.NoError(t, err)
		log.Debugf("user %s pair created", nil, user.name)
	}
	for _, user := range usrs {
		p, err := New(user.passphrase)
		require.NoError(t, err)
		privPath := fmt.Sprintf("./openpgp/%s.gpg", user.name)
		priv, err := os.ReadFile(privPath)
		require.NoError(t, err)
		p.AddPrivateKey(priv)
		for _, u := range usrs {
			pubPath := fmt.Sprintf("./openpgp/%s.gpg.pub", u.name)
			pub, err := os.ReadFile(pubPath)
			require.NoError(t, err)
			p.AddPublicKey(pub)
		}
		testname := fmt.Sprintf("test encryption with %s", user.name)
		t.Run(testname, func(t *testing.T) {
			encrypted, err := p.EncryptWithKeys(msg)
			require.NoError(t, err)
			encryptedText, err := encrypted.GetArmored()
			require.NoError(t, err)
			log.Debugf("Encrypted message: \n%s", nil, encryptedText)
			for _, uu := range usrs {
				pp, err := New(uu.passphrase)
				require.NoError(t, err)
				privPath := fmt.Sprintf("./openpgp/%s.gpg", uu.name)
				priv, err := os.ReadFile(privPath)
				require.NoError(t, err)
				pp.AddPrivateKey(priv)
				for _, uuu := range usrs {
					pubPath := fmt.Sprintf("./openpgp/%s.gpg.pub", uuu.name)
					pub, err := os.ReadFile(pubPath)
					require.NoError(t, err)
					pp.AddPublicKey(pub)
				}
				subtestname := fmt.Sprintf("test decryption with %s", uu.name)
				t.Run(subtestname, func(t *testing.T) {
					got, err := pp.DecryptWithKey(encrypted.GetBinary())
					require.NoError(t, err)
					require.Equal(t, msg, got)
					log.Debug(nil, string(got))
				})
			}

		})
	}
	content, err := os.ReadDir("./openpgp")
	require.NoError(t, err)
	for _, file := range content {
		path := fmt.Sprintf("./openpgp/%s", file.Name())
		err := os.Remove(path)
		require.NoError(t, err)
		log.Debugf("%s deleted", nil, path)
	}
	err = os.RemoveAll("./openpgp")
	require.NoError(t, err)
	log.Debug(nil, "./openpgp")
}

// func TestStream(t *testing.T) {
// 	p, err := New("passphrase")
// 	require.NoError(t, err)
// 	err = p.GeneratePair("test-name", "test@email.com")
// 	require.NoError(t, err)
// 	msg := []byte("testz<fdbafrbarfbadsfbvzfdbzdnbzd v<sdrfv<s")
// 	t.Run("compressed encryption", func(t *testing.T) {
// 		f, err := p.NewPGPStream("test.gpg")
// 		require.NoError(t, err)
// 		fmt.Fprint(f.ReadWriter(), string(msg))
// 		err = f.Close()
// 		require.NoError(t, err)
// 	})
// 	t.Run("compressed decryption", func(t *testing.T) {
// 		f, err := p.NewPGPStream("test.gpg")
// 		require.NoError(t, err)
// 		t.Log(f.buf.Bytes())
// 		data, err := io.ReadAll(f.ReadWriter())
// 		require.NoError(t, err)
// 		require.Equal(t, msg, data)
// 		err = f.Close()
// 		require.NoError(t, err)
// 	})

// }

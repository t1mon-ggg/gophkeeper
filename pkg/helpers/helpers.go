package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"syscall"

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

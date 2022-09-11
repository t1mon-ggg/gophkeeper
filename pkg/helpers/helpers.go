package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
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

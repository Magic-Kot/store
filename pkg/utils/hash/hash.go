package hash

import (
	"fmt"

	"crypto/sha256"
)

const salt = "qwertyuiop0sdfghjklzxcvbnm1234567890QWE8TYUIOPASDFGHJKLZXCVBNM"

// GenerateHash - генерация хеша
func GenerateHash(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

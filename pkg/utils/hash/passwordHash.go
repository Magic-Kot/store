package hash

import (
	"fmt"

	"crypto/sha256"
)

const salt = "qwertyuiop0sdfghjklzxcvbnm1234567890QWE8TYUIOPASDFGHJKLZXCVBNM"

// GeneratePasswordHash - генерация хеша пароля
func GeneratePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

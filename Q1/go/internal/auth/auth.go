package auth

import (
	"crypto/sha512"
	"fmt"
)

func HashSHA512(input string) string {
	hash := sha512.Sum512([]byte(input))
	return fmt.Sprintf("%x", hash)
}

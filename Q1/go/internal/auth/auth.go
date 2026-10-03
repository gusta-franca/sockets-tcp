/*
   Data de criação: 30/09/2026
   Estudante: Gustavo Martins França
   Suporte à função hash SHA512
*/

package auth

import (
	"crypto/sha512"
	"fmt"
)

// Helper to hash passwords
func HashSHA512(input string) string {
	hash := sha512.Sum512([]byte(input))
	return fmt.Sprintf("%x", hash)
}

package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)

	rand.Read(token)

	hexToken := hex.EncodeToString(token)
	return hexToken, nil
}

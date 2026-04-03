package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

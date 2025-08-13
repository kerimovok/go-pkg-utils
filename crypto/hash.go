package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashSHA256Hex returns the hex-encoded SHA-256 hash of the input string.
func HashSHA256Hex(input string) string {
	sum := sha256.Sum256([]byte(input))
	return hex.EncodeToString(sum[:])
}

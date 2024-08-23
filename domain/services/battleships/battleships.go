package battleships

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Commit generates a SHA-256 hash commitment
// value: the value to commit to (e.g., 0 for water, 1 for a ship)
// randomness: a random value to ensure the commitment is hiding
func Commit(value int, randomness int) string {
	// Combine the value and randomness into a single string
	data := fmt.Sprintf("%d%d", value, randomness)

	// Create a new SHA-256 hash
	hash := sha256.New()

	// Write the data to the hash
	hash.Write([]byte(data))

	// Calculate the SHA-256 checksum and return it as a hexadecimal string
	return hex.EncodeToString(hash.Sum(nil))
}

func VerifyProof(miss bool, randomValue int, commitment string) bool {
	var computedCommitment string
	if miss {
		computedCommitment = Commit(0, randomValue)
	} else {
		computedCommitment = Commit(1, randomValue)
	}
	return computedCommitment == commitment
}

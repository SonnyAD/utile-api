package battleships

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Commit(t *testing.T) {
	// Example usage
	value := 1          // 1 represents a ship
	randomness := 12345 // A random value

	commitment := Commit(value, randomness)
	t.Log("Commitment:", commitment)
	assert.Equal(t, "03c28c828cc2b2558d975399118363de6dbf96a7ac82dfa53621c524319349a1", commitment)
}

package api

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CalculatePi(t *testing.T) {
	tt := map[string]struct {
		value    *big.Float
		expected string
	}{
		"pi": {
			value:    chudnovsky(10000),
			expected: fmt.Sprint(math.Pi),
		},
		"tau": {
			value:    chudnovskyTau(10000),
			expected: fmt.Sprint(math.Pi * 2.0),
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			strVal := fmt.Sprintf("%.10000f", tc.value)
			assert.True(t, strings.HasPrefix(strVal, tc.expected))
		})
	}
}

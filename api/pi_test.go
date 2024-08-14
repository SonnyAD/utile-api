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

func Benchmark_CalculatePi(b *testing.B) {
	tt := map[string]struct {
		value int
	}{
		"10": {
			value: 10,
		},
		"100": {
			value: 100,
		},
		"1000": {
			value: 1000,
		},
		"10000": {
			value: 10000,
		},
		"100000": {
			value: 100000,
		},
		"1000000": {
			value: 1000000,
		},
	}

	for name, tc := range tt {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				chudnovsky(tc.value)
			}
		})
	}
}

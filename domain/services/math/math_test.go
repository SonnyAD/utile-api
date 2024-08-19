package math

import (
	"bufio"
	"fmt"
	"math"
	"math/big"
	"os"
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
			value:    Chudnovsky(10000),
			expected: fmt.Sprint(math.Pi),
		},
		"tau": {
			value:    ChudnovskyTau(10000),
			expected: fmt.Sprint(math.Pi * 2.0),
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			strVal := fmt.Sprintf("%.1000000f", tc.value)

			f, err := os.Create("/tmp/" + name)
			assert.NoError(t, err)

			defer f.Close()

			w := bufio.NewWriter(f)
			_, err = w.WriteString(strVal)
			assert.NoError(t, err)

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
				Chudnovsky(tc.value)
			}
		})
	}
}

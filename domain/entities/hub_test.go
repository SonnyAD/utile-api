package entities

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_CountOnlinePlayers(t *testing.T) {
	t.Parallel()

	client := NewClient(nil, nil)
	client2 := NewClient(nil, nil)

	tt := map[string]struct {
		givenClients  []*Client
		expectedCount int
	}{
		"none": {
			givenClients:  make([]*Client, 0),
			expectedCount: 0,
		},
		"one": {
			givenClients:  []*Client{client},
			expectedCount: 1,
		},
		"one-too": {
			givenClients:  []*Client{client, client},
			expectedCount: 1,
		},
		"two": {
			givenClients:  []*Client{client, client2},
			expectedCount: 2,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			hub := NewHub()
			go hub.Run(ctx)

			for _, client := range tc.givenClients {
				hub.Register <- client
				time.Sleep(10 * time.Millisecond)
			}

			assert.Equal(t, tc.expectedCount, hub.CountOnlinePlayers())
		})
	}
}

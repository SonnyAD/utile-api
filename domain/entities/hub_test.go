package entities

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_CountOnlinePlayers(t *testing.T) {
	t.Parallel()

	client := NewClient(nil, nil)
	client2 := NewClient(nil, nil)
	uuid1 := uuid.NewString()
	uuid2 := uuid.NewString()

	tt := map[string]struct {
		givenClients   []*Client
		givenPlayerIDs []string
		expectedCount  int
	}{
		"none": {
			givenClients:  make([]*Client, 0),
			expectedCount: 0,
		},
		"one": {
			givenClients:   []*Client{client},
			givenPlayerIDs: []string{uuid1},
			expectedCount:  1,
		},
		"one-too": {
			givenClients:   []*Client{client, client},
			givenPlayerIDs: []string{uuid1, uuid1},
			expectedCount:  1,
		},
		"two": {
			givenClients:   []*Client{client, client2},
			givenPlayerIDs: []string{uuid1, uuid2},
			expectedCount:  2,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			hub := NewHub()
			go hub.Run(ctx)

			for i, client := range tc.givenClients {
				hub.Register <- client
				time.Sleep(10 * time.Millisecond)
				client.SetPlayerID(tc.givenPlayerIDs[i])
				hub.RecordPlayer(client.PlayerID, client)
			}

			assert.Equal(t, tc.expectedCount, hub.CountOnlinePlayers())
		})
	}
}

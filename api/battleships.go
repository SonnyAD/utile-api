package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	battleships "utile.space/api/domain/entities"
	"utile.space/api/utils"
)

var (
	hub *battleships.Hub
)

func getHub(ctx context.Context) *battleships.Hub {
	if hub == nil {
		hub = battleships.NewHub()
		go hub.Run(ctx)
	}
	return hub
}

var upgraderBattleShips = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// @Summary		BattleshipsWebsocket to play battleships with another player or computer
// @Description	Websocket to open to play battleships
// @Tags			battleships
// @Success		101
// @Router			/math/ws [get]
func BattleshipsWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgraderBattleShips.Upgrade(w, r, nil)
	if err != nil {
		fmt.Print("upgrade:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer c.Close()

	hub := getHub(context.Background())
	client := battleships.NewClient(hub, c)
	client.Hub.Register <- client

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		client.WritePump()
		wg.Done()
	}()

	go func() {
		client.ReadPump()
		wg.Done()
	}()

	wg.Wait()
}

// @Summary		BattleshipsStats to get stats on the multiplayer state of the game
// @Description	To get stats on multiplayer state of the battleships game
// @Tags			battleships
// @Success		200 {object}	StatsResult
// @Router			/math/stats [get]
func BattleshipsStats(w http.ResponseWriter, r *http.Request) {
	var stats StatsResult
	hub := getHub(context.Background())
	stats.OnlinePlayersCount = hub.CountOnlinePlayers()
	stats.PendingMatchesCount = hub.CountPendingMatches()
	stats.OngoingMatchesCount = hub.CountOngoingMatches()
	stats.FinishedMatchesCount = hub.CountFinishedMatches()
	stats.TotalMatchesCount = hub.CountTotalMatches()

	utils.Output(w, r.Header["Accept"], stats, strconv.Itoa(stats.OnlinePlayersCount))
}

type StatsResult struct {
	XMLName              xml.Name `json:"-" xml:"stats" yaml:"-"`
	OnlinePlayersCount   int      `json:"onlinePlayers" xml:"OnlinePlayers" yaml:"onlinePlayers"`
	PendingMatchesCount  int      `json:"pendingMatches" xml:"PendingMatches" yaml:"pendingMatches"`
	OngoingMatchesCount  int      `json:"ongoingMatches" xml:"OngoingMatches" yaml:"ongoingMatches"`
	FinishedMatchesCount int      `json:"finishedMatches" xml:"FinishedMatches" yaml:"finishedMatches"`
	TotalMatchesCount    int      `json:"totalMatches" xml:"TotalMatches" yaml:"totalMatches"`
}

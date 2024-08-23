package api

import (
	"fmt"
	"net/http"
	"time"

	battleships "utile.space/api/domain/entities"
)

var (
	hub *battleships.Hub
)

func getHub() *battleships.Hub {
	if hub == nil {
		hub = battleships.NewHub()
		go hub.Run()
	}
	return hub
}

// @Summary		BattleshipsWebsocket to play battleships with another player or computer
// @Description	Websocket to open to play battleships
// @Tags			battleships
// @Success		101
// @Router			/math/ws [get]
func BattleshipsWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Print("upgrade:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer c.Close()

	client := battleships.NewClient(getHub(), c)
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
	for {
		time.Sleep(time.Second)
	}
}

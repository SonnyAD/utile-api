package api

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/websocket"
	"utile.space/api/domain/services/math"
	"utile.space/api/utils"
)

// @Summary		Pi Value
// @Description	Calculate Pi value up to 10K decimals
// @Tags			math
// @Produce		json,xml,application/yaml,plain
// @Success		200	{object}	BigNumberResult
// @Router			/math/pi [get]
func CalculatePi(w http.ResponseWriter, r *http.Request) {
	pi := math.Chudnovsky(10000)

	var answer BigNumberResult
	answer.Name = "Pi"
	answer.Value = fmt.Sprintf("%.10000f", pi)

	utils.Output(w, r.Header["Accept"], answer, answer.Value)
}

// @Summary		Tau Value
// @Description	Calculate Tau value up to 10K decimals
// @Tags			math
// @Produce		json,xml,application/yaml,plain
// @Success		200	{object}	BigNumberResult
// @Router			/math/tau [get]
func CalculateTau(w http.ResponseWriter, r *http.Request) {
	tau := math.ChudnovskyTau(10000)

	var answer BigNumberResult
	answer.Name = "Tau"
	answer.Value = fmt.Sprintf("%.10000f", tau)

	utils.Output(w, r.Header["Accept"], answer, answer.Value)
}

type BigNumberResult struct {
	XMLName xml.Name `json:"-" xml:"bignumber" yaml:"-"`
	Name    string   `json:"name" xml:"name" yaml:"name"`
	Value   string   `json:"value" xml:"value" yaml:"value"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// @Summary		MathWebsocket to get pi and tau by page up to 1M digits
// @Description	Websocket to get pi and tau by page up to 1M digits. It will switch protocols as requested.
// @Tags			math
// @Success		101
// @Router			/math/ws [get]
func MathWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Print("upgrade:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			continue
		}
		log.Printf("recv: %s", message)

		r := regexp.MustCompile(`^(pi|tau)\s+([0-9]+),\s*([0-9]+)$`)
		subMatch := r.FindStringSubmatch(string(message))

		// pi or tau
		if subMatch != nil {
			page, err := strconv.Atoi(subMatch[2])
			if err != nil {
				log.Println("write:", err)
				continue
			}
			pageSize, err := strconv.Atoi(subMatch[3])
			if err != nil {
				log.Println("write:", err)
				continue
			}

			err = c.WriteMessage(mt, []byte(math.ReadNextPage(subMatch[1], page, pageSize)))
			if err != nil {
				log.Println("write:", err)
				continue
			}
		}
	}
}

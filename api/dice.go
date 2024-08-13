package api

import (
	"encoding/xml"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"utile.space/api/utils"
)

// @Summary		Roll a dice
// @Description	Endpoint to roll a dice of the given number of faces
// @Tags			dice
// @Produce		json,xml,application/yaml,plain
// @Param			dice	path		int	true	"Number of faces of the dice between 2 and 100"
// @Success		200		{object}	DieResult
// @Router			/d{dice} [get]
func RollDice(w http.ResponseWriter, r *http.Request) {
	dice, err := strconv.Atoi(mux.Vars(r)["dice"])

	if err != nil {
		http.Error(w, "Die not found", http.StatusNotFound)
		return
	}

	var roll DieResult
	roll.Die = dice
	roll.Result = rand.Intn(dice) + 1

	utils.Output(w, r.Header["Accept"], roll, strconv.Itoa(roll.Result))
}

type DieResult struct {
	XMLName xml.Name `json:"-" xml:"dieresult" yaml:"-"`
	Die     int      `json:"die" xml:"die" yaml:"die"`
	Result  int      `json:"result" xml:"result" yaml:"result"`
}

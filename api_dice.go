package main

import (
	"encoding/xml"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func RollDice(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)

	dice, err := strconv.Atoi(mux.Vars(r)["dice"])

	if err != nil {
		http.Error(w, "Die not found", http.StatusNotFound)
		return
	}

	var roll DieResult
	roll.Die = dice
	roll.Result = rand.Intn(dice) + 1

	output(w, r.Header["Accept"], roll, strconv.Itoa(roll.Result))
}

type DieResult struct {
	XMLName xml.Name `json:"-" xml:"dieresult" yaml:"-"`
	Die     int      `json:"die" xml:"die" yaml:"die"`
	Result  int      `json:"result" xml:"result" yaml:"result"`
}

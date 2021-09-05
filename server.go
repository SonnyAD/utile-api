package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

func EmptyResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)

	var health Health
	health.Status = "up"

	output(w, r.Header["Accept"], health, health.Status)
}

type Health struct {
	XMLName xml.Name `json:"-" xml:"health" yaml:"-"`
	Version string   `json:"version,omitempty" xml:"version,omitempty" yaml:"version,omitempty"`
	Status  string   `json:"status" xml:"status" yaml:"status"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func output(w http.ResponseWriter, accept []string, v interface{}, plain string) {
	if contains(accept, "application/json") {
		reply, _ := json.Marshal(v)
		fmt.Fprintf(w, string(reply))
	} else if contains(accept, "application/xml") {
		reply, _ := xml.Marshal(v)
		fmt.Fprintf(w, string(reply))
	} else if contains(accept, "application/yaml") {
		reply, _ := yaml.Marshal(v)
		fmt.Fprintf(w, string(reply))
	} else {
		fmt.Fprintf(w, plain)
	}
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", EmptyResponse).Methods("GET")

	// NOTE: need to use non capturing group with (?:pattern) below because capturing group are not supported
	router.HandleFunc("/d{dice:(?:100|1[0-9]|[2-9][0-9]?)}", RollDice).Methods("GET")
	router.HandleFunc("/dns/{domain}", DNSResolve).Methods("GET")
	router.HandleFunc("/dns/mx/{domain}", MXResolve).Methods("GET")
	router.HandleFunc("/dns/cname/{domain}", CNAMEResolve).Methods("GET")
	router.HandleFunc("/dns/txt/{domain}", TXTResolve).Methods("GET")
	router.HandleFunc("/dns/ns/{domain}", NSResolve).Methods("GET")
	router.HandleFunc("/dns/caa/{domain}", CAAResolve).Methods("GET")

	router.HandleFunc("/status", HealthCheck).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", router))
}

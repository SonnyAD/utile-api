package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

func RollDice(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)

	dice, err := strconv.Atoi(mux.Vars(r)["dice"])

	if err != nil {
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, strconv.Itoa(rand.Intn(dice)+1))
}

func DNSResolve(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)

	domain := mux.Vars(r)["domain"]

	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, "1.1.1.1:53")
		},
	}

	var dns DNSResolved

	ip, _ := resolver.LookupHost(context.Background(), domain)

	dns.Addresses = ip

	if contains(r.Header["Accept"], "text/html") {
		fmt.Fprintf(w, ip[0])
	} else if contains(r.Header["Accept"], "application/json") {
		reply, _ := json.Marshal(dns)
		fmt.Fprintf(w, string(reply))
	} else if contains(r.Header["Accept"], "application/xml") {
		reply, _ := xml.Marshal(dns)
		fmt.Fprintf(w, string(reply))
	} else if contains(r.Header["Accept"], "application/yaml") {
		reply, _ := yaml.Marshal(dns)
		fmt.Fprintf(w, string(reply))
	} else {
		fmt.Fprintf(w, ip[0])
	}
}

type DNSResolved struct {
	Addresses []string `json:"addresses" xml:"addresses" yaml:"addresses"`
}

func DoHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
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

	router.HandleFunc("/", DoHealthCheck).Methods("GET")

	// NOTE: need to use non capturing group with (?:pattern) below because capturing group are not supported
	router.HandleFunc("/d{dice:(?:100|1[0-9]|[2-9][0-9]?)}", RollDice).Methods("GET")
	router.HandleFunc("/dns/{domain}", DNSResolve).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", router))
}

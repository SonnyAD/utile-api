package utils

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

func EnableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	})
}

func Output(w http.ResponseWriter, accept []string, v interface{}, plain string) {
	if contains(accept, "application/json") {
		reply, _ := json.Marshal(v)
		fmt.Fprint(w, string(reply))
	} else if contains(accept, "application/xml") {
		reply, _ := xml.Marshal(v)
		fmt.Fprint(w, string(reply))
	} else if contains(accept, "application/yaml") {
		reply, _ := yaml.Marshal(v)
		fmt.Fprint(w, string(reply))
	} else {
		fmt.Fprint(w, plain)
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

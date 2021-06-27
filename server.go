package main

import (
  "fmt"
  "net/http"
  "math/rand"
  "strconv"
  "regexp"
)

func home(w http.ResponseWriter, r *http.Request) {

  // evaluate only if i in [2-100]
  re := regexp.MustCompile("/d(1[0-9]|[2-9][0-9]?|100)")
  match := re.FindStringSubmatch(r.URL.Path)

  if match == nil {
    http.Error(w, "Bad request", http.StatusBadRequest)
    return
  }

  dice, err := strconv.Atoi(match[1])

  if err != nil {
    http.Error(w, "Unknown error", http.http.StatusInternalError)
    return
  }

  fmt.Fprintf(w, strconv.Itoa(rand.Intn(dice)+1))
}

func main() {
    http.HandleFunc("/", home)
    http.ListenAndServe(":3000", nil)
}


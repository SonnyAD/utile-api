package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"utile.space/api/api"
	_ "utile.space/api/docs"
	"utile.space/api/utils"
)

func EmptyResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// @Summary		Healthcheck
// @Description	Get the status of the API
// @Tags			health
// @Produce		json,xml,application/yaml,plain
// @Success		200	{object}	Health
// @Router			/status [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {

	utils.EnableCors(&w)

	var health Health
	health.Status = "up"

	version, present := os.LookupEnv("API_VERSION")
	if present {
		health.Version = version
	}

	utils.Output(w, r.Header["Accept"], health, health.Status)
}

type Health struct {
	XMLName xml.Name `json:"-" xml:"health" yaml:"-"`
	Version string   `json:"version,omitempty" xml:"version,omitempty" yaml:"version,omitempty"`
	Status  string   `json:"status" xml:"status" yaml:"status"`
}

// @title			utile.space Open API
// @version		1.0
// @description	The collection of free API from utile.space, the Swiss Army Knife webtool.
//
// @contact.name	API Support
// @contact.email	api@utile.space
//
// @license.name	utile.space API License
// @license.url	https://utile.space/api/
//
// @BasePath		/api
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", EmptyResponse).Methods(http.MethodGet)

	// NOTE: need to use non capturing group with (?:pattern) below because capturing group are not supported
	router.HandleFunc("/d{dice:(?:100|1[0-9]|[2-9][0-9]?)}", api.RollDice).Methods(http.MethodGet)
	router.HandleFunc("/dns/{domain}", api.DNSResolve).Methods(http.MethodGet)
	router.HandleFunc("/dns/mx/{domain}", api.MXResolve).Methods(http.MethodGet)
	router.HandleFunc("/dns/cname/{domain}", api.CNAMEResolve).Methods(http.MethodGet)
	router.HandleFunc("/dns/txt/{domain}", api.TXTResolve).Methods(http.MethodGet)
	router.HandleFunc("/dns/ns/{domain}", api.NSResolve).Methods(http.MethodGet)
	router.HandleFunc("/dns/caa/{domain}", api.CAAResolve).Methods(http.MethodGet)
	router.HandleFunc("/dns/aaaa/{domain}", api.AAAAResolve).Methods(http.MethodGet)
	router.HandleFunc("/dns/dmarc/{domain}", api.DMARCResolve).Methods(http.MethodGet)
	router.HandleFunc("/dns/ptr/{ip}", api.PTRResolve).Methods(http.MethodGet)
	router.HandleFunc("/links", api.GetLinksPage).Methods(http.MethodGet)

	router.HandleFunc("/status", HealthCheck).Methods(http.MethodGet)

	router.PathPrefix("/docs/").Handler(httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	port, present := os.LookupEnv("PORT")
	if !present {
		port = "3000"
	}

	log.Fatal(http.ListenAndServe(":"+port, router))
}

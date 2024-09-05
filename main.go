package main

import (
	"encoding/xml"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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

func initLogging() {
	log.SetLevel(log.DebugLevel)
	//log.SetFormatter(&log.JSONFormatter{})
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
	initLogging()

	router := mux.NewRouter()

	router.Use(utils.EnableCors)

	apiRouter := router.PathPrefix("/api").Subrouter()

	router.HandleFunc("/", EmptyResponse).Methods(http.MethodGet)
	apiRouter.HandleFunc("/", EmptyResponse).Methods(http.MethodGet)

	// NOTE: need to use non capturing group with (?:pattern) below because capturing group are not supported
	apiRouter.HandleFunc("/d{dice:(?:100|1[0-9]|[2-9][0-9]?)}", api.RollDice).Methods(http.MethodGet)
	apiRouter.HandleFunc("/dns/{domain}", api.DNSResolve).Methods(http.MethodGet)
	apiRouter.HandleFunc("/dns/mx/{domain}", api.MXResolve).Methods(http.MethodGet)
	apiRouter.HandleFunc("/dns/cname/{domain}", api.CNAMEResolve).Methods(http.MethodGet)
	apiRouter.HandleFunc("/dns/txt/{domain}", api.TXTResolve).Methods(http.MethodGet)
	apiRouter.HandleFunc("/dns/ns/{domain}", api.NSResolve).Methods(http.MethodGet)
	apiRouter.HandleFunc("/dns/caa/{domain}", api.CAAResolve).Methods(http.MethodGet)
	apiRouter.HandleFunc("/dns/aaaa/{domain}", api.AAAAResolve).Methods(http.MethodGet)
	apiRouter.HandleFunc("/dns/dmarc/{domain}", api.DMARCResolve).Methods(http.MethodGet)
	apiRouter.HandleFunc("/dns/ptr/{ip}", api.PTRResolve).Methods(http.MethodGet)
	apiRouter.HandleFunc("/links", api.GetLinksPage).Methods(http.MethodGet)
	apiRouter.HandleFunc("/math/pi", api.CalculatePi).Methods(http.MethodGet)
	apiRouter.HandleFunc("/math/tau", api.CalculateTau).Methods(http.MethodGet)
	apiRouter.HandleFunc("/math/ws", api.MathWebsocket).Methods(http.MethodGet)
	apiRouter.HandleFunc("/battleships/ws", api.BattleshipsWebsocket).Methods(http.MethodGet)
	apiRouter.HandleFunc("/battleships/stats", api.BattleshipsStats).Methods(http.MethodGet)

	apiRouter.HandleFunc("/status", HealthCheck).Methods(http.MethodGet)

	apiRouter.PathPrefix("/docs/").Handler(httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	port, present := os.LookupEnv("PORT")
	if !present {
		port = "3000"
	}

	log.Info("Starting server on port ", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

package main

import (
	"os"

	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/polarbroadband/goto/util"
)

var (
	// container name
	HOST = "HELLO"
)

func init() {
	// config package level default logger
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
}

type SVR struct {
	*util.API
}

func main() {
	api := SVR{
		API: &util.API{
			Log: log.WithField("owner", HOST),
		},
	}

	// config http server
	r := mux.NewRouter()
	rru := r.PathPrefix("").Subrouter()
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{`*`}),
		//handlers.AllowedHeaders([]string{"content-type", "X-Csrf-Token", "withcredentials", "credentials",}),
		handlers.AllowedHeaders([]string{"content-type"}),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "HEAD"}),
	)

	rru.HandleFunc("/", api.healtz).Methods("GET")

	api.Log.Fatal(http.ListenAndServe(":80", cors(rru)))
}

// healtz response k8s health check probe
func (a *SVR) healtz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_e := util.NewExeErr("healtz", HOST, "rest_API")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "done", "greeting": HOST}); err != nil {
		a.Error(w, 500, _e.String("erroneous api response", err))
	}
}

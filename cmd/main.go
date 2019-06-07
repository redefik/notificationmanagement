package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/redefik/notificationmanagement/config"
	"github.com/redefik/notificationmanagement/repository"
	"github.com/redefik/notificationmanagement/resthandler"
	"log"
	"net/http"
)

// healthCheck exposed the endpoint used to check the state of the microservice
func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

var configurationFile = flag.String("config", "config/config.json", "Location of the config file.")

func main() {
	flag.Parse()
	// Load the configuration parameters of the microservice
	err := config.SetConfiguration(*configurationFile)
	if err != nil {
		log.Panicln(err)
	}
	repository.InitializeDynamoDbClient()
	r := mux.NewRouter()
	// Register the handlers for the various HTTP requests
	r.HandleFunc("/", healthCheck).Methods(http.MethodGet)
	r.HandleFunc("/notification_management/api/v1.0/course", resthandler.NewCourse).Methods(http.MethodPost)
	r.HandleFunc("/notification_management/api/v1.0/course", resthandler.DeleteCourse).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(config.Configuration.ListeningAddress, r))
}

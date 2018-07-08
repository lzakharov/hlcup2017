package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// configuration := LoadConfiguration("config.json")
	// host, db := configuration.DB.Host, configuration.DB.Database

	router := mux.NewRouter()

	router.HandleFunc("/users/{id:[0-9]+}", nil).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}/visits", nil).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}", nil).Methods("POST")
	router.HandleFunc("/users/new", nil).Methods("POST")

	router.HandleFunc("/locations/{id:[0-9]+}", nil).Methods("GET")
	router.HandleFunc("/locations/{id:[0-9]+}/avg", nil).Methods("GET")
	router.HandleFunc("/locations/{id:[0-9]+}", nil).Methods("POST")
	router.HandleFunc("/locations/new", nil).Methods("POST")

	router.HandleFunc("/visits/{id:[0-9]+}", nil).Methods("GET")
	router.HandleFunc("/visits/{id:[0-9]+}", nil).Methods("POST")
	router.HandleFunc("/visits/new", nil).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lzakharov/hlcup2017/config"
	"github.com/lzakharov/hlcup2017/handlers"
	"github.com/lzakharov/hlcup2017/models"
	"github.com/lzakharov/hlcup2017/utils"
)

func DumbHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, I'm Web Server!"))
}

func main() {
	c := config.LoadConfiguration("config.json")

	db := c.DB
	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.Name, db.SSLMode)
	models.InitDatabase(db.Driver, dataSourceName)
	models.CreateSchema(db.Schema)

	utils.LoadData(c.Data)

	r := mux.NewRouter()

	r.HandleFunc("/", DumbHandler).Methods("GET")

	r.HandleFunc("/users/{id:[0-9]+}", handlers.GetUser).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}/visits", handlers.GetUserVisits).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", nil).Methods("POST")
	r.HandleFunc("/users/new", handlers.CreateUser).Methods("POST")

	r.HandleFunc("/locations/{id:[0-9]+}", handlers.GetLocation).Methods("GET")
	r.HandleFunc("/locations/{id:[0-9]+}/avg", handlers.GetLocationAverageMark).Methods("GET")
	r.HandleFunc("/locations/{id:[0-9]+}", nil).Methods("POST")
	r.HandleFunc("/locations/new", handlers.CreateLocation).Methods("POST")

	r.HandleFunc("/visits/{id:[0-9]+}", handlers.GetVisit).Methods("GET")
	r.HandleFunc("/visits/{id:[0-9]+}", nil).Methods("POST")
	r.HandleFunc("/visits/new", handlers.CreateVisit).Methods("POST")

	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	log.Fatal(http.ListenAndServe(addr, r))
}

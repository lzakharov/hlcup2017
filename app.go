package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"log"
	"net/http"
)

var emptyJSON = map[string]string{}

type App struct {
	Router   *mux.Router
	Database *Database
}

func (a *App) Initialize(c *Config) error {
	a.Database = new(Database)
	if err := a.Database.Initialize(c.DBConfig); err != nil {
		return err
	}
	if err := LoadData(c.Data, a.Database); err != nil {
		return nil
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()

	return nil
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/users/{id:[0-9]+}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/users/{id:[0-9]+}/visits", a.getUserVisits).Methods("GET")
	a.Router.HandleFunc("/users/{id:[0-9]+}", a.updateUser).Methods("POST")
	a.Router.HandleFunc("/users/new", a.createUser).Methods("POST")

	a.Router.HandleFunc("/locations/{id:[0-9]+}", a.getLocation).Methods("GET")
	a.Router.HandleFunc("/locations/{id:[0-9]+}/avg", a.getLocationAverageMark).Methods("GET")
	a.Router.HandleFunc("/locations/{id:[0-9]+}", a.updateLocation).Methods("POST")
	a.Router.HandleFunc("/locations/new", a.createLocation).Methods("POST")

	a.Router.HandleFunc("/visits/{id:[0-9]+}", a.getVisit).Methods("GET")
	a.Router.HandleFunc("/visits/{id:[0-9]+}", a.updateVisit).Methods("POST")
	a.Router.HandleFunc("/visits/new", a.createVisit).Methods("POST")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := vars["id"]

	user, err := a.Database.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) getUserVisits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := vars["id"]

	filter := new(PlaceFilter)
	decoder := schema.NewDecoder()
	if err := decoder.Decode(filter, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	places, err := a.Database.GetUserVisits(id, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(places); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.Database.InsertUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(emptyJSON); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := vars["id"]

	user := new(User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.Database.UpdateUser(id, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(emptyJSON); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) getLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := vars["id"]

	location, err := a.Database.GetLocation(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(location); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) getLocationAverageMark(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := vars["id"]

	filter := new(LocationFilter)
	decoder := schema.NewDecoder()
	if err := decoder.Decode(filter, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	avg, err := a.Database.GetLocationAverageMark(id, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(avg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) createLocation(w http.ResponseWriter, r *http.Request) {
	location := new(Location)
	if err := json.NewDecoder(r.Body).Decode(location); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.Database.InsertLocation(location); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(emptyJSON); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) updateLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := vars["id"]

	location := new(Location)
	if err := json.NewDecoder(r.Body).Decode(location); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.Database.UpdateLocation(id, location); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(emptyJSON); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) getVisit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := vars["id"]

	visit, err := a.Database.GetVisit(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(visit); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) createVisit(w http.ResponseWriter, r *http.Request) {
	visit := new(Visit)
	if err := json.NewDecoder(r.Body).Decode(visit); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.Database.InsertVisit(visit); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(emptyJSON); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) updateVisit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := vars["id"]

	visit := new(Visit)
	if err := json.NewDecoder(r.Body).Decode(visit); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.Database.UpdateVisit(id, visit); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(emptyJSON); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

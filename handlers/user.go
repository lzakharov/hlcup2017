package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/lzakharov/hlcup2017/models"
)

// GetUser returns specified user.
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := vars["id"]

	user, err := models.GetUser(id)
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

// GetUserVisits returns visits of specified user.
func GetUserVisits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := vars["id"]

	filter := new(models.PlaceFilter)
	decoder := schema.NewDecoder()
	if err := decoder.Decode(filter, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	places, err := models.GetUserVisits(id, filter)
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

// CreateUser creates new user.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := models.InsertUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{})
}

// UpdateUser updates specified user.
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := vars["id"]

	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := models.UpdateUser(id, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(EmptyJSON)
}

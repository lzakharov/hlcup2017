package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lzakharov/hlcup2017/models"
)

// GetLocation returns specified location.
func GetLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := vars["id"]

	location, err := models.GetLocation(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(location); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

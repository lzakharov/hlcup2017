package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lzakharov/hlcup2017/models"
)

// GetVisit returns specified visit.
func GetVisit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := vars["id"]

	visit, err := models.GetVisit(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(visit); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

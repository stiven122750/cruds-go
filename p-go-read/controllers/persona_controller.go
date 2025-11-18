package controllers

import (
	"encoding/json"
	"net/http"
	services "p-go-read/service"

	"github.com/gorilla/mux"
)

func ObtenerPersonaHandler(w http.ResponseWriter, r *http.Request) {
	documento := mux.Vars(r)["documento"]

	persona, err := services.ObtenerPersona(documento)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(persona)
}

func ObtenerTodasLasPersonasHandler(w http.ResponseWriter, r *http.Request) {
	personas, err := services.ObtenerTodasLasPersonas()
	if err != nil {
		http.Error(w, "Error consultando personas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(personas)
}

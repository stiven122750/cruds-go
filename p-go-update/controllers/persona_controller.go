package controllers

import (
	"encoding/json"
	"net/http"
	"p-go-update/models"
	services "p-go-update/service"

	"github.com/gorilla/mux"
)

func ActualizarPersonaHandler(w http.ResponseWriter, r *http.Request) {
	documento := mux.Vars(r)["documento"]

	var data models.PersonaUpdate
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if err := services.ActualizarPersona(documento, data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

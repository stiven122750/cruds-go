package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/stiven122750/cruds-go/p-go-delete/services"
)

func EliminarPersona(w http.ResponseWriter, r *http.Request) {
	// Obtener el documento desde los parámetros de la URL
	vars := mux.Vars(r)
	documento := vars["documento"]

	if documento == "" {
		http.Error(w, "El documento es obligatorio", http.StatusBadRequest)
		return
	}

	// Llamar al servicio
	err := services.EliminarPersona(documento)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Persona eliminada exitosamente"})
}

package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stiven122750/cruds-go/p-go-delete/models"
	"github.com/stiven122750/cruds-go/p-go-delete/services"
)

// ================================================================
//
//	MOCK DEL REPOSITORIO
//
// ================================================================
type MockRepo struct {
	Existe bool
	Error  bool
}

func (m MockRepo) ObtenerPersonaPorDocumento(documento string) (models.Persona, error) {
	if m.Error {
		return models.Persona{}, errors.New("error")
	}
	if !m.Existe {
		return models.Persona{}, errors.New("no existe")
	}
	return models.Persona{Documento: documento}, nil
}

func (m MockRepo) EliminarPersonaPorDocumento(documento string) error {
	if m.Error {
		return errors.New("error eliminando")
	}
	return nil
}

// ================================================================
//
//	UNIT TESTS
//
// ================================================================
func TestEliminarPersona_OK(t *testing.T) {
	// Inyectar mock: persona sí existe
	services.SetPersonaRepository(MockRepo{Existe: true})

	req := httptest.NewRequest("DELETE", "/eliminar-persona/123", nil)
	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/eliminar-persona/{documento}", EliminarPersona).Methods("DELETE")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Se esperaba 200, se obtuvo %d", w.Code)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["mensaje"] != "Persona eliminada exitosamente" {
		t.Fatalf("Mensaje incorrecto: %v", resp)
	}
}

func TestEliminarPersona_NoExiste(t *testing.T) {
	// Inyectar mock: persona NO existe
	services.SetPersonaRepository(MockRepo{Existe: false})

	req := httptest.NewRequest("DELETE", "/eliminar-persona/999", nil)
	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/eliminar-persona/{documento}", EliminarPersona).Methods("DELETE")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Se esperaba 404, se obtuvo %d", w.Code)
	}
}

func TestEliminarPersona_BadRequest(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/eliminar-persona/", nil)
	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/eliminar-persona/{documento}", EliminarPersona).Methods("DELETE")

	router.ServeHTTP(w, req)

	// gorilla devuelve 404 cuando falta un parámetro
	if w.Code != http.StatusNotFound {
		t.Fatalf("Se esperaba 404, se obtuvo %d", w.Code)
	}
}

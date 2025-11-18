package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	services "p-go-read/service"
	"testing"

	"p-go-read/models"

	"github.com/gorilla/mux"
)

// ----------------- Mock del REPO -----------------

type MockRepo struct {
	PersonaFake models.Persona
	ErrorFake   bool
}

func (m MockRepo) ObtenerPersonaPorDocumento(doc string) (models.Persona, error) {
	if m.ErrorFake {
		return models.Persona{}, errors.New("no encontrado")
	}
	return m.PersonaFake, nil
}

func (m MockRepo) ObtenerTodasLasPersonas() ([]models.Persona, error) {
	if m.ErrorFake {
		return nil, errors.New("error")
	}
	return []models.Persona{m.PersonaFake}, nil
}

// Inyectar mock al service
func setMockRepo() {
	services.SetPersonaRepository(MockRepo{
		PersonaFake: models.Persona{
			Documento: "123",
			Nombre:    "Steven",
		},
	})
}

// ----------------- TESTS -----------------

func TestObtenerPersonaHandler_OK(t *testing.T) {
	setMockRepo()

	req, _ := http.NewRequest("GET", "/personas/123", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/personas/{documento}", ObtenerPersonaHandler)

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperaba 200, obtuve %d", rr.Code)
	}

	var p models.Persona
	json.Unmarshal(rr.Body.Bytes(), &p)

	if p.Documento != "123" {
		t.Fatalf("Documento incorrecto: %s", p.Documento)
	}
}

func TestObtenerPersonaHandler_NoExiste(t *testing.T) {
	services.SetPersonaRepository(MockRepo{
		ErrorFake: true,
	})

	req, _ := http.NewRequest("GET", "/personas/999", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/personas/{documento}", ObtenerPersonaHandler)

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("Esperaba 404, obtuve %d", rr.Code)
	}
}

func TestObtenerTodasLasPersonasHandler_OK(t *testing.T) {
	setMockRepo()

	req, _ := http.NewRequest("GET", "/personas", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/personas", ObtenerTodasLasPersonasHandler)

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperaba 200, obtuve %d", rr.Code)
	}

	var list []models.Persona
	json.Unmarshal(rr.Body.Bytes(), &list)

	if len(list) != 1 {
		t.Fatalf("Esperaba 1 persona, obtuve %d", len(list))
	}
}

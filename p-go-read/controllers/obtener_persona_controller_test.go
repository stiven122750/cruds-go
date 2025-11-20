package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"p-go-read/models"
	"p-go-read/service"
	"testing"

	"github.com/gorilla/mux"
)

// ===============================
// MOCK REPO PARA PRUEBAS UNITARIAS
// ===============================

type MockRepoUnit struct {
	FakePersona models.Persona
	ShouldError bool
}

func (m MockRepoUnit) ObtenerPersonaPorDocumento(documento string) (models.Persona, error) {
	if m.ShouldError {
		return models.Persona{}, errors.New("persona no encontrada")
	}
	return m.FakePersona, nil
}

func (m MockRepoUnit) ObtenerTodasLasPersonas() ([]models.Persona, error) {
	if m.ShouldError {
		return nil, errors.New("error en BD")
	}
	return []models.Persona{m.FakePersona}, nil
}

func resetRepoUnit() {
	service.SetPersonaRepository(nil)
}

// ===============================
// PRUEBAS UNITARIAS DEL CONTROLLER
// ===============================

func TestObtenerPersonaHandler_OK(t *testing.T) {
	defer resetRepoUnit()

	fake := models.Persona{
		Documento: "123",
		Nombre:    "Steven",
	}
	service.SetPersonaRepository(MockRepoUnit{FakePersona: fake})

	req, _ := http.NewRequest("GET", "/personas/123", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/personas/{documento}", ObtenerPersonaHandler)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtenido %d", rr.Code)
	}

	var resp models.Persona
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp.Documento != "123" {
		t.Fatalf("documento incorrecto")
	}
}

func TestObtenerPersonaHandler_NoExiste(t *testing.T) {
	defer resetRepoUnit()

	service.SetPersonaRepository(MockRepoUnit{ShouldError: true})

	req, _ := http.NewRequest("GET", "/personas/999", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/personas/{documento}", ObtenerPersonaHandler)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("esperado 404")
	}
}

func TestObtenerPersonaHandler_DocumentoVacio(t *testing.T) {
	defer resetRepoUnit()

	service.SetPersonaRepository(MockRepoUnit{})

	req, _ := http.NewRequest("GET", "/personas/%20", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/personas/{documento}", ObtenerPersonaHandler)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("esperado 404 por documento vacío")
	}
}

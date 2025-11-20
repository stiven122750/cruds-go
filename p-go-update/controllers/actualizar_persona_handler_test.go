package controllers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"p-go-update/models"
	services "p-go-update/service"
	"testing"

	"github.com/gorilla/mux"
)

//
// ===========================
//          MOCKS
// ===========================
//

// Mock del ReadClient
type MockReadClient struct {
	Persona models.Persona
	Error   error
}

func (m MockReadClient) ObtenerPersona(documento string) (models.Persona, error) {
	return m.Persona, m.Error
}

// Mock del repositorio
type MockRepo struct {
	UpdateErr error
}

func (m MockRepo) ActualizarPersonaPorDocumento(documento string, data models.PersonaUpdate) error {
	return m.UpdateErr
}

//
// ===========================
//           TESTS
// ===========================
//

func TestActualizarPersonaHandler_OK(t *testing.T) {

	// Mock para ReadClient
	services.SetReadClient(MockReadClient{
		Persona: models.Persona{Documento: "123"},
		Error:   nil,
	})

	// Mock para Repo
	services.SetPersonaRepository(MockRepo{
		UpdateErr: nil,
	})

	body := `{"nombre":"Nuevo","edad":20}`
	req := httptest.NewRequest("PUT", "/personas/123", bytes.NewBuffer([]byte(body)))
	req = mux.SetURLVars(req, map[string]string{"documento": "123"})

	rr := httptest.NewRecorder()

	ActualizarPersonaHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("esperado 200, obtenido %d", rr.Code)
	}

	expected := `{"status":"ok"}`
	if rr.Body.String() != expected {
		t.Fatalf("esperado %s, obtenido %s", expected, rr.Body.String())
	}
}

func TestActualizarPersonaHandler_JSONInvalido(t *testing.T) {

	req := httptest.NewRequest("PUT", "/personas/123", bytes.NewBuffer([]byte("{invalid json}")))
	req = mux.SetURLVars(req, map[string]string{"documento": "123"})

	rr := httptest.NewRecorder()

	ActualizarPersonaHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("esperado 400, obtenido %d", rr.Code)
	}
}

func TestActualizarPersonaHandler_ErrorEnUpdate(t *testing.T) {

	services.SetReadClient(MockReadClient{
		Persona: models.Persona{Documento: "555"},
		Error:   nil,
	})

	services.SetPersonaRepository(MockRepo{
		UpdateErr: errors.New("falló actualización"),
	})

	body := `{"nombre":"Nuevo"}`
	req := httptest.NewRequest("PUT", "/personas/555", bytes.NewBuffer([]byte(body)))
	req = mux.SetURLVars(req, map[string]string{"documento": "555"})

	rr := httptest.NewRecorder()

	ActualizarPersonaHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("esperado 400, obtenido %d", rr.Code)
	}
}

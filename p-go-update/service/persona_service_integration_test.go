package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"p-go-update/models"
	"testing"
)

// ==============================
// Mock del ReadClient
// ==============================

type MockReadClient struct {
	URL         string
	ShouldError bool
}

func (m MockReadClient) ObtenerPersona(documento string) (models.Persona, error) {

	if m.ShouldError {
		return models.Persona{}, errors.New("error forzado desde mock")
	}

	// Consumir la URL simulada
	resp, err := http.Get(m.URL + "/personas/" + documento)
	if err != nil {
		return models.Persona{}, err
	}
	defer resp.Body.Close()

	var p models.Persona
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return models.Persona{}, err
	}

	return p, nil
}

// ==============================
// TEST: ObtenerPersonaEnRead
// ==============================

func TestObtenerPersonaEnRead_Exitoso(t *testing.T) {

	// Servidor HTTP falso que simula el microservicio READ
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Validar la ruta consultada
		if r.URL.Path != "/personas/123" {
			t.Fatalf("Ruta inesperada: %s", r.URL.Path)
		}

		resp := models.Persona{
			Documento: "123",
			Nombre:    "Maria",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Inyectar cliente falso en el servicio
	SetReadClient(MockReadClient{URL: server.URL})

	// Ejecutar la llamada
	p, err := ReadClientInstance.ObtenerPersona("123")
	if err != nil {
		t.Fatalf("No esperaba error, obtuve: %v", err)
	}

	// Validaciones
	if p.Documento != "123" {
		t.Errorf("Documento incorrecto, esperado 123, recibido %s", p.Documento)
	}

	if p.Nombre != "Maria" {
		t.Errorf("Nombre incorrecto, esperado Maria, recibido %s", p.Nombre)
	}
}

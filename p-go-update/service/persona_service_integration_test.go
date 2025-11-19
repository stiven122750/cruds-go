package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"p-go-update/models"
	"testing"
)

// Mock del ReadClient para usar el servidor falso
type MockReadClient struct {
	URL string
}

func (m MockReadClient) ObtenerPersona(documento string) (models.Persona, error) {
	resp, err := http.Get(m.URL)
	if err != nil {
		return models.Persona{}, err
	}
	defer resp.Body.Close()

	var p models.Persona
	json.NewDecoder(resp.Body).Decode(&p)
	return p, nil
}

func TestObtenerPersonaEnRead_Exitoso(t *testing.T) {

	// Servidor HTTP falso simulando READ
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := models.Persona{Documento: "123", Nombre: "Maria"}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Inyectar cliente falso
	SetReadClient(MockReadClient{URL: server.URL})

	// Ejecutar solicitud simulada
	p, err := ReadClientInstance.ObtenerPersona("123")
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}

	if p.Nombre != "Maria" {
		t.Errorf("Nombre incorrecto, esperado Maria, recibido %s", p.Nombre)
	}
}

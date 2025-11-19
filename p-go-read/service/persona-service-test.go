package services

import (
	"errors"
	"p-go-read/models"
	"testing"
)

// ===============================
// MOCK REAL DEL REPOSITORIO
// ===============================

type MockPersonaRepo struct {
	FakePersona models.Persona
	ShouldError bool
}

func (m MockPersonaRepo) ObtenerPersonaPorDocumento(doc string) (models.Persona, error) {
	if m.ShouldError {
		return models.Persona{}, errors.New("persona no encontrada")
	}
	return m.FakePersona, nil
}

func (m MockPersonaRepo) ObtenerTodasLasPersonas() ([]models.Persona, error) {
	if m.ShouldError {
		return nil, errors.New("error al consultar")
	}
	return []models.Persona{m.FakePersona}, nil
}

// ===============================
// TEST: ObtenerPersona()
// ===============================

func TestObtenerPersona_OK(t *testing.T) {

	personaFake := models.Persona{
		Documento: "123",
		Nombre:    "Steven",
	}

	SetPersonaRepository(MockPersonaRepo{
		FakePersona: personaFake,
	})

	p, err := ObtenerPersona("123")

	if err != nil {
		t.Fatalf("No esperaba error, obtuve: %v", err)
	}

	if p.Documento != "123" {
		t.Fatalf("Documento incorrecto: %s", p.Documento)
	}
}

func TestObtenerPersona_Vacio(t *testing.T) {
	SetPersonaRepository(MockPersonaRepo{})

	_, err := ObtenerPersona("")

	if err == nil {
		t.Fatalf("Esperaba error por documento vacío")
	}
}

func TestObtenerPersona_NoExiste(t *testing.T) {
	SetPersonaRepository(MockPersonaRepo{
		ShouldError: true,
	})

	_, err := ObtenerPersona("999")

	if err == nil {
		t.Fatalf("Esperaba error porque la persona no existe")
	}
}

// ===============================
// TEST: ObtenerTodasLasPersonas()
// ===============================

func TestObtenerTodas_OK(t *testing.T) {

	personaFake := models.Persona{
		Documento: "1",
		Nombre:    "Test",
	}

	SetPersonaRepository(MockPersonaRepo{
		FakePersona: personaFake,
	})

	list, err := ObtenerTodasLasPersonas()

	if err != nil {
		t.Fatalf("No esperaba error: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("Esperaba 1 persona, recibí %d", len(list))
	}
}

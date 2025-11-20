package service

import (
	"errors"
	"p-go-update/models"
	"testing"
)

//
// ====================
//      MOCKS
// ====================
//

// Mock del ReadClient exclusivo de este test
type MockReadClient_Update struct {
	Resp models.Persona
	Err  error
}

func (m MockReadClient_Update) ObtenerPersona(documento string) (models.Persona, error) {
	return m.Resp, m.Err
}

// Mock del repositorio exclusivo de este test
type MockRepo_Update struct {
	Err error
}

func (m MockRepo_Update) ActualizarPersonaPorDocumento(doc string, data models.PersonaUpdate) error {
	return m.Err
}

//
// =============================
//        TEST UNITARIOS
// =============================
//

func TestActualizarPersona_Success(t *testing.T) {

	SetReadClient(MockReadClient_Update{
		Resp: models.Persona{Documento: "123"},
		Err:  nil,
	})

	SetPersonaRepository(MockRepo_Update{Err: nil})

	update := models.PersonaUpdate{
		Nombre:   "Juan",
		Apellido: "Lopez",
		Edad:     30,
	}

	err := ActualizarPersona("123", update)
	if err != nil {
		t.Errorf("Esperaba nil, obtuvo: %v", err)
	}
}

func TestActualizarPersona_DocumentoVacio(t *testing.T) {

	update := models.PersonaUpdate{
		Nombre: "Juan",
	}

	err := ActualizarPersona("", update)

	if err == nil || err.Error() != "el documento no puede estar vacío" {
		t.Errorf("Error esperado 'el documento no puede estar vacío', obtuvo: %v", err)
	}
}

func TestActualizarPersona_UpdateInvalido(t *testing.T) {

	update := models.PersonaUpdate{} // vacío → inválido

	err := ActualizarPersona("123", update)

	if err == nil || err.Error() != "no hay campos para actualizar" {
		t.Errorf("Error esperado 'no hay campos para actualizar', obtuvo: %v", err)
	}
}

func TestActualizarPersona_ReadClientNoConfigurado(t *testing.T) {

	SetReadClient(nil) // ← PROVOCAMOS EL ERROR
	SetPersonaRepository(MockRepo_Update{})

	update := models.PersonaUpdate{
		Nombre: "Nuevo",
	}

	err := ActualizarPersona("123", update)

	if err == nil || err.Error() != "ReadClient no configurado" {
		t.Errorf("Error esperado 'ReadClient no configurado', obtuvo: %v", err)
	}
}

func TestActualizarPersona_ReadClientFalla(t *testing.T) {

	SetReadClient(MockReadClient_Update{
		Resp: models.Persona{},
		Err:  errors.New("la persona no existe"),
	})

	SetPersonaRepository(MockRepo_Update{})

	update := models.PersonaUpdate{
		Nombre: "Nuevo",
	}

	err := ActualizarPersona("123", update)

	if err == nil || err.Error() != "la persona no existe" {
		t.Errorf("Error esperado 'la persona no existe', obtuvo: %v", err)
	}
}

func TestActualizarPersona_RepoFalla(t *testing.T) {

	SetReadClient(MockReadClient_Update{
		Resp: models.Persona{Documento: "123"},
		Err:  nil,
	})

	SetPersonaRepository(MockRepo_Update{
		Err: errors.New("fallo al actualizar"),
	})

	update := models.PersonaUpdate{
		Nombre: "Nuevo",
	}

	err := ActualizarPersona("123", update)

	if err == nil || err.Error() != "fallo al actualizar" {
		t.Errorf("Error esperado 'fallo al actualizar', obtuvo: %v", err)
	}
}

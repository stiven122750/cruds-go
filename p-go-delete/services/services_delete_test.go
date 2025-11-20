package services

import (
	"errors"
	"testing"

	"github.com/stiven122750/cruds-go/p-go-delete/models"
	"github.com/stiven122750/cruds-go/p-go-delete/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestEliminarPersonaExitosa(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	Repo = mockRepo

	documento := "123"
	persona := models.Persona{Documento: documento}

	// 1️⃣ Simular que la persona existe
	mockRepo.
		On("ObtenerPersonaPorDocumento", documento).
		Return(persona, nil)

	// 2️⃣ Simular eliminación exitosa
	mockRepo.
		On("EliminarPersonaPorDocumento", documento).
		Return(nil)

	err := EliminarPersona(documento)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestEliminarPersonaDocumentoVacio(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	Repo = mockRepo

	err := EliminarPersona("")

	assert.EqualError(t, err, "el documento no puede estar vacío")
}

func TestEliminarPersonaNoExiste(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	Repo = mockRepo

	documento := "999"

	// 1️⃣ Simular que la persona NO existe
	mockRepo.
		On("ObtenerPersonaPorDocumento", documento).
		Return(models.Persona{}, errors.New("not found"))

	err := EliminarPersona(documento)

	// 2️⃣ El service traduce el error
	assert.EqualError(t, err, "no existe persona con ese documento")

	mockRepo.AssertExpectations(t)
}

package services_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stiven122750/p-go-delete/services"
	"github.com/stiven122750/p-go-delete/tests/mocks"
	"github.com/stiven122750/p-go-delete/models"
)

func TestEliminarPersonaExitosa(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	services.Repo = mockRepo

	documento := "123"
	persona := models.Persona{Documento: documento}

	// Simular que la persona existe
	mockRepo.On("ObtenerPersonaPorDocumento", documento).Return(persona, nil)
	// Simular eliminación exitosa
	mockRepo.On("EliminarPersonaPorDocumento", documento).Return(nil)

	err := services.EliminarPersona(documento)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestEliminarPersonaDocumentoVacio(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	services.Repo = mockRepo

	err := services.EliminarPersona("")
	assert.EqualError(t, err, "el documento no puede estar vacío")
}

func TestEliminarPersonaNoExiste(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	services.Repo = mockRepo

	documento := "999"

	// Simular que la persona NO existe
	mockRepo.On("ObtenerPersonaPorDocumento", documento).Return(models.Persona{}, errors.New("not found"))

	err := services.EliminarPersona(documento)

	assert.EqualError(t, err, "no existe persona con ese documento")
}

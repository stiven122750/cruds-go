package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stiven122750/cruds-go/p-go-delete/services"
	"github.com/stiven122750/cruds-go/p-go-delete/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestEliminarPersonaController_Success(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	services.SetPersonaRepository(mockRepo)

	documento := "123"

	// Mock flujo exitoso: persona existe y se elimina correctamente
	mockRepo.On("ObtenerPersonaPorDocumento", documento).Return(nil)
	mockRepo.On("EliminarPersonaPorDocumento", documento).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/eliminar-persona/"+documento, nil)
	rr := httptest.NewRecorder()

	controllers.EliminarPersona(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Persona eliminada exitosamente")
	mockRepo.AssertExpectations(t)
}

func TestEliminarPersonaController_NoExiste(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	services.SetPersonaRepository(mockRepo)

	documento := "999"
	mockRepo.On("ObtenerPersonaPorDocumento", documento).Return(errors.New("not found"))

	req := httptest.NewRequest(http.MethodDelete, "/eliminar-persona/"+documento, nil)
	rr := httptest.NewRecorder()

	controllers.EliminarPersona(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "no existe persona con ese documento")
	mockRepo.AssertExpectations(t)
}

func TestEliminarPersonaController_DocumentoVacio(t *testing.T) {
	mockRepo := new(mocks.MockPersonaRepo)
	services.SetPersonaRepository(mockRepo)

	req := httptest.NewRequest(http.MethodDelete, "/eliminar-persona/", nil)
	rr := httptest.NewRecorder()

	controllers.EliminarPersona(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "el documento no puede estar vacío")
}

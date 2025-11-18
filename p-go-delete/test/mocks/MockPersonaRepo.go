package mocks

import (
	"github.com/stiven122750/p-go-delete/models"
	"github.com/stretchr/testify/mock"
)

type MockPersonaRepo struct {
	mock.Mock
}

func (m *MockPersonaRepo) ObtenerPersonaPorDocumento(doc string) (models.Persona, error) {
	args := m.Called(doc)
	return args.Get(0).(models.Persona), args.Error(1)
}

func (m *MockPersonaRepo) EliminarPersonaPorDocumento(doc string) error {
	args := m.Called(doc)
	return args.Error(0)
}

package mocks

import (
	"p-go-read/models"

	"github.com/stretchr/testify/mock"
)

type MockPersonaRepo struct {
	mock.Mock
}

// -----------------------------
// READ: obtener 1 persona
// -----------------------------
func (m *MockPersonaRepo) ObtenerPersonaPorDocumento(doc string) (models.Persona, error) {
	args := m.Called(doc)

	// Si el valor es nil, devolvemos el zero-value
	persona, _ := args.Get(0).(models.Persona)

	return persona, args.Error(1)
}

// -----------------------------
// READ: obtener todas
// -----------------------------
func (m *MockPersonaRepo) ObtenerTodasLasPersonas() ([]models.Persona, error) {
	args := m.Called()

	// Si el valor es nil, devolvemos slice vacío
	list, ok := args.Get(0).([]models.Persona)
	if !ok {
		list = []models.Persona{}
	}

	return list, args.Error(1)
}

// -----------------------------
// DELETE
// -----------------------------
func (m *MockPersonaRepo) EliminarPersonaPorDocumento(doc string) error {
	args := m.Called(doc)
	return args.Error(0)
}

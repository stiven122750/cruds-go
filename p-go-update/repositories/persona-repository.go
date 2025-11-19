package repositories

import "p-go-update/models"

// Interface para permitir mocking en pruebas
type PersonaRepository interface {
	ActualizarPersonaPorDocumento(documento string, data models.PersonaUpdate) error
}

package repositories

import "github.com/danysoftdev/p-go-create/models"

// Interface para permitir mocking en pruebas
type PersonaRepository interface {
	EliminarPersonaPorDocumento(documento string) error
	ObtenerPersonaPorDocumento(documento string) (models.Persona, error)
}

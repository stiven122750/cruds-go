package repositories

import "github.com/stiven122750/p-go-delete/models"

// Interface para permitir mocking en pruebas
type PersonaRepository interface {
	EliminarPersonaPorDocumento(documento string) error
	ObtenerPersonaPorDocumento(documento string) (models.Persona, error)
}

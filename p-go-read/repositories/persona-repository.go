package repositories

import "p-go-read/models"

// Interface para permitir mocking en pruebas
type PersonaRepository interface {
	ObtenerPersonaPorDocumento(documento string) (models.Persona, error)
	ObtenerTodasLasPersonas() ([]models.Persona, error)
}

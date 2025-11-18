package services

import (
	"errors"
	"p-go-read/models"
	"p-go-read/repositories"
	"strings"
)

var Repo repositories.PersonaRepository

func SetPersonaRepository(r repositories.PersonaRepository) {
	Repo = r
}

// Validar documento antes de consultar
func ValidarDocumento(documento string) error {
	if strings.TrimSpace(documento) == "" {
		return errors.New("el documento no puede estar vacío")
	}
	return nil
}

// Servicio: obtener una persona por documento
func ObtenerPersona(documento string) (models.Persona, error) {
	if err := ValidarDocumento(documento); err != nil {
		return models.Persona{}, err
	}

	return Repo.ObtenerPersonaPorDocumento(documento)
}

// Servicio: obtener todas las personas
func ObtenerTodasLasPersonas() ([]models.Persona, error) {
	return Repo.ObtenerTodasLasPersonas()
}

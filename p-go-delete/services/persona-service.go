package services

import (
	"errors"
	"strings"

	"github.com/stiven122750/cruds-go/p-go-delete/models"
	"github.com/stiven122750/cruds-go/p-go-delete/repositories"
)

var Repo repositories.PersonaRepository

func SetPersonaRepository(r repositories.PersonaRepository) {
	Repo = r
}

// Validar documento antes de eliminar
func ValidarDocumento(documento string) error {
	if strings.TrimSpace(documento) == "" {
		return errors.New("el documento no puede estar vacío")
	}
	return nil
}

// Servicio: eliminar persona
func EliminarPersona(documento string) error {
	if err := ValidarDocumento(documento); err != nil {
		return err
	}

	// Verificar si existe
	_, err := Repo.ObtenerPersonaPorDocumento(documento)
	if err != nil {
		return errors.New("no existe persona con ese documento")
	}

	// Proceder a eliminar
	return Repo.EliminarPersonaPorDocumento(documento)
}

// Este se usa también en GET o validaciones opcionales
func ObtenerPersona(documento string) (models.Persona, error) {
	return Repo.ObtenerPersonaPorDocumento(documento)
}

package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"p-go-update/models"
	"p-go-update/repositories"
	"strings"
	"time"
)

var Repo repositories.PersonaRepository

func SetPersonaRepository(r repositories.PersonaRepository) {
	Repo = r
}

// ---------------- VALIDACIONES ----------------

func ValidarDocumento(documento string) error {
	if strings.TrimSpace(documento) == "" {
		return errors.New("el documento no puede estar vacío")
	}
	return nil
}

func ValidarUpdate(data models.PersonaUpdate) error {
	if data.Nombre == "" && data.Apellido == "" && data.Edad == 0 {
		return errors.New("no hay campos para actualizar")
	}
	return nil
}

// ---------------- CONSULTAR AL MICROSERVICIO READ ----------------

func ObtenerPersonaEnRead(documento string) (models.Persona, error) {

	// ⚠️ Usar el host del contenedor en Docker Compose
	url := fmt.Sprintf("http://read-service:5000/personas/%s", documento)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return models.Persona{}, errors.New("error comunicándose con el microservicio READ")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return models.Persona{}, errors.New("la persona no existe")
	}

	if resp.StatusCode == http.StatusOK {
		var p models.Persona
		if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
			return models.Persona{}, errors.New("respuesta inválida del read-service")
		}
		return p, nil
	}

	return models.Persona{}, fmt.Errorf("read-service retornó código inesperado: %d", resp.StatusCode)
}

// ---------------- SERVICE UPDATE ----------------

func ActualizarPersona(documento string, data models.PersonaUpdate) error {

	// Validar documento
	if err := ValidarDocumento(documento); err != nil {
		return err
	}

	// Validar body
	if err := ValidarUpdate(data); err != nil {
		return err
	}

	// Consultar existencia con microservicio READ
	_, err := ObtenerPersonaEnRead(documento)
	if err != nil {
		return err
	}

	// Actualizar la persona en el repositorio
	return Repo.ActualizarPersonaPorDocumento(documento, data)
}

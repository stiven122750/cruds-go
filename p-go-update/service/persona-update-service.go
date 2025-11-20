package service

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

//
// =========================
//   REPOSITORIO (MOCKEABLE)
// =========================
//

var Repo repositories.PersonaRepository

func SetPersonaRepository(r repositories.PersonaRepository) {
	Repo = r
}

//
// ===================================
//   INTERFAZ PARA EL CLIENTE DE READ
// ===================================
//

// Interfaz para permitir mocks en tests
type ReadClient interface {
	ObtenerPersona(documento string) (models.Persona, error)
}

var ReadClientInstance ReadClient

func SetReadClient(c ReadClient) {
	ReadClientInstance = c
}

//
// ==========================================
//   IMPLEMENTACIÓN REAL DEL CLIENTE DE READ
// ==========================================
//

type ReadServiceClient struct{}

func (r ReadServiceClient) ObtenerPersona(documento string) (models.Persona, error) {

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

//
// ==================
//     VALIDACIONES
// ==================
//

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

//
// ======================
//      SERVICE UPDATE
// ======================
//

func ActualizarPersona(documento string, data models.PersonaUpdate) error {

	// Validar documento
	if err := ValidarDocumento(documento); err != nil {
		return err
	}

	// Validar body
	if err := ValidarUpdate(data); err != nil {
		return err
	}

	// Consultar existencia en microservicio READ (mockeable)
	if ReadClientInstance == nil {
		return errors.New("ReadClient no configurado")
	}

	_, err := ReadClientInstance.ObtenerPersona(documento)
	if err != nil {
		return err
	}

	// Actualizar en repositorio
	return Repo.ActualizarPersonaPorDocumento(documento, data)
}

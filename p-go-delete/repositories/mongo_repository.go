package repositories

import (
	"context"
	"time"

	"github.com/stiven122750/cruds-go/p-go-delete/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

// Recibimos la colección desde main (inyección de dependencias)
func SetCollection(c *mongo.Collection) {
	collection = c
}

// Buscar persona (opcional para validación antes del delete)
func ObtenerPersonaPorDocumento(documento string) (models.Persona, error) {
	var persona models.Persona
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"documento": documento}).Decode(&persona)
	return persona, err
}

// Eliminar persona
func EliminarPersonaPorDocumento(documento string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"documento": documento})
	return err
}

// Implementación real del repositorio
type RealPersonaRepository struct{}

func (r RealPersonaRepository) EliminarPersonaPorDocumento(doc string) error {
	return EliminarPersonaPorDocumento(doc)
}

func (r RealPersonaRepository) ObtenerPersonaPorDocumento(doc string) (models.Persona, error) {
	return ObtenerPersonaPorDocumento(doc)
}

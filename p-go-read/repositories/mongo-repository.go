package repositories

import (
	"context"
	"p-go-read/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

// Inyección de dependencias desde main
func SetCollection(c *mongo.Collection) {
	collection = c
}

// Obtener una persona por documento
func ObtenerPersonaPorDocumento(documento string) (models.Persona, error) {
	var persona models.Persona
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"documento": documento}).Decode(&persona)
	return persona, err
}

// Obtener todas las personas
func ObtenerTodasLasPersonas() ([]models.Persona, error) {
	var personas []models.Persona
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var persona models.Persona
		if err := cursor.Decode(&persona); err != nil {
			return nil, err
		}
		personas = append(personas, persona)
	}

	return personas, cursor.Err()
}

// Implementación real del repositorio
type RealPersonaRepository struct{}

// READ individual
func (r RealPersonaRepository) ObtenerPersonaPorDocumento(doc string) (models.Persona, error) {
	return ObtenerPersonaPorDocumento(doc)
}

// READ lista
func (r RealPersonaRepository) ObtenerTodasLasPersonas() ([]models.Persona, error) {
	return ObtenerTodasLasPersonas()
}

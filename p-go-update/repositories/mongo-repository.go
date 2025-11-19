package repositories

import (
	"context"
	"p-go-update/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

// Inyección de dependencias desde main
func SetCollection(c *mongo.Collection) {
	collection = c
}

// Actualizar persona por documento
func ActualizarPersonaPorDocumento(documento string, data models.PersonaUpdate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filtro := bson.M{"documento": documento}
	update := bson.M{"$set": data}

	result, err := collection.UpdateOne(ctx, filtro, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// Implementación real del repositorio
type RealPersonaRepository struct{}

// UPDATE
func (r RealPersonaRepository) ActualizarPersonaPorDocumento(documento string, data models.PersonaUpdate) error {
	return ActualizarPersonaPorDocumento(documento, data)
}

package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Collection *mongo.Collection
var client *mongo.Client

func ConectarMongo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")
	collectionName := os.Getenv("COLLECTION_NAME")

	if uri == "" || dbName == "" || collectionName == "" {
		return fmt.Errorf("faltan variables de entorno: MONGO_URI, MONGO_DB o COLLECTION_NAME")
	}

	clientOptions := options.Client().ApplyURI(uri)
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Confirmar conexión
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	Collection = client.Database(dbName).Collection(collectionName)
	log.Println("✅ Conectado a MongoDB correctamente.")
	return nil
}

func CerrarMongo() error {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return client.Disconnect(ctx)
	}
	return nil
}

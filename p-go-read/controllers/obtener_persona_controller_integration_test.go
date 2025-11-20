package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"p-go-read/models"
	"p-go-read/service"
)

const mongoURI = "mongodb://localhost:27017"

type testMongoRepo struct {
	collection *mongo.Collection
}

// ====================
// Implementación inline
// ====================

func (r *testMongoRepo) ObtenerPersonaPorDocumento(documento string) (models.Persona, error) {
	var p models.Persona
	err := r.collection.FindOne(context.Background(), bson.M{"documento": documento}).Decode(&p)
	if err == mongo.ErrNoDocuments {
		return models.Persona{}, errors.New("persona no encontrada")
	}
	return p, err
}

func (r *testMongoRepo) ObtenerTodasLasPersonas() ([]models.Persona, error) {
	cur, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	var personas []models.Persona
	for cur.Next(context.Background()) {
		var p models.Persona
		cur.Decode(&p)
		personas = append(personas, p)
	}
	return personas, nil
}

// ====================
// Setup Mongo
// ====================

func setupMongo(t *testing.T) (*mongo.Client, *mongo.Collection) {
	t.Helper()

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("Error creando cliente Mongo: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		t.Fatalf("Error conectando a Mongo: %v", err)
	}

	col := client.Database("test_db").Collection("personas")
	col.DeleteMany(context.Background(), bson.D{})

	return client, col
}

func TestObtenerPersonaHandler_Integration(t *testing.T) {
	client, col := setupMongo(t)
	defer client.Disconnect(context.Background())

	// Inyectar repositorio inline
	service.SetPersonaRepository(&testMongoRepo{collection: col})

	// Insertar datos
	_, err := col.InsertOne(context.Background(), models.Persona{
		Documento: "123",
		Nombre:    "Juan",
	})
	if err != nil {
		t.Fatalf("Error insertando datos: %v", err)
	}

	req := httptest.NewRequest("GET", "/personas/123", nil)
	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/personas/{documento}", ObtenerPersonaHandler)

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Status esperado 200, recibido %d", rr.Code)
	}

	var resp models.Persona
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp.Documento != "123" {
		t.Fatalf("Documento esperado '123', recibido '%s'", resp.Documento)
	}
}

func TestObtenerTodasLasPersonasHandler_Integration(t *testing.T) {
	client, col := setupMongo(t)
	defer client.Disconnect(context.Background())

	service.SetPersonaRepository(&testMongoRepo{collection: col})

	col.InsertOne(context.Background(), models.Persona{Documento: "1"})
	col.InsertOne(context.Background(), models.Persona{Documento: "2"})

	req := httptest.NewRequest("GET", "/personas", nil)
	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/personas", ObtenerTodasLasPersonasHandler)

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Status esperado 200, recibido %d", rr.Code)
	}

	var resp []models.Persona
	json.Unmarshal(rr.Body.Bytes(), &resp)

	if len(resp) != 2 {
		t.Fatalf("Esperado 2 personas, recibido %d", len(resp))
	}
}

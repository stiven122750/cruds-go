package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"p-go-read/models"
	"p-go-read/repositories"
	services "p-go-read/service"
	"testing"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestMongo(t *testing.T) *mongo.Collection {
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Error conectando Mongo: %v", err)
	}

	db := client.Database("test_read_controller_db")
	col := db.Collection("personas")

	col.DeleteMany(context.Background(), bson.M{})

	repositories.SetCollection(col)
	services.SetPersonaRepository(repositories.RealPersonaRepository{})

	return col
}

func TestIntegration_Controller_ObtenerPersona(t *testing.T) {
	col := setupTestMongo(t)

	col.InsertOne(context.Background(), models.Persona{
		Documento: "333",
		Nombre:    "Integration Test",
	})

	router := mux.NewRouter()
	router.HandleFunc("/personas/{documento}", ObtenerPersonaHandler)

	req, _ := http.NewRequest("GET", "/personas/333", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperaba 200, obtuve %d", rr.Code)
	}

	var p models.Persona
	json.Unmarshal(rr.Body.Bytes(), &p)

	if p.Documento != "333" {
		t.Fatalf("Documento incorrecto: %s", p.Documento)
	}
}

func TestIntegration_Controller_ObtenerTodas(t *testing.T) {
	col := setupTestMongo(t)

	col.InsertMany(context.Background(), []interface{}{
		models.Persona{Documento: "1", Nombre: "A"},
		models.Persona{Documento: "2", Nombre: "B"},
	})

	router := mux.NewRouter()
	router.HandleFunc("/personas", ObtenerTodasLasPersonasHandler)

	req, _ := http.NewRequest("GET", "/personas", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Esperaba 200, obtuve %d", rr.Code)
	}

	var list []models.Persona
	json.Unmarshal(rr.Body.Bytes(), &list)

	if len(list) != 2 {
		t.Fatalf("Esperaba 2, obtuve %d", len(list))
	}
}

package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	services "p-go-update/service"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"p-go-update/models"
)

const mongoURI = "mongodb://localhost:27017"

// ===============================
//
//	Mock REAL del ReadClient
//
// ===============================
type readClientMock struct{}

func (c readClientMock) ObtenerPersona(documento string) (models.Persona, error) {
	// Siempre retorna "persona existe"
	return models.Persona{
		Documento: documento,
		Nombre:    "Mock",
		Apellido:  "User",
		Edad:      99,
	}, nil
}

// ===============================
//
//	Repositorio de Mongo inline
//
// ===============================
type testMongoRepo struct {
	col *mongo.Collection
}

func (r *testMongoRepo) ActualizarPersonaPorDocumento(documento string, data models.PersonaUpdate) error {
	update := bson.M{"$set": bson.M{}}

	if data.Nombre != "" {
		update["$set"].(bson.M)["nombre"] = data.Nombre
	}
	if data.Apellido != "" {
		update["$set"].(bson.M)["apellido"] = data.Apellido
	}
	if data.Edad != 0 {
		update["$set"].(bson.M)["edad"] = data.Edad
	}

	_, err := r.col.UpdateOne(context.Background(), bson.M{"documento": documento}, update)
	return err
}

// ===============================
//
//	Setup Mongo para tests
//
// ===============================
func setupMongo(t *testing.T) (*mongo.Client, *mongo.Collection) {
	t.Helper()

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("Error creando cliente Mongo: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		t.Fatalf("Error conectando a Mongo: %v", err)
	}

	col := client.Database("test_db").Collection("personas_update")
	col.DeleteMany(context.Background(), bson.D{})

	return client, col
}

// ===============================
//  TEST DE INTEGRACIÓN
// ===============================

func TestActualizarPersonaHandler_Integration(t *testing.T) {
	client, col := setupMongo(t)
	defer client.Disconnect(context.Background())

	// Insertamos datos iniciales
	_, err := col.InsertOne(context.Background(), bson.M{
		"documento": "123",
		"nombre":    "Juan",
		"apellido":  "Lopez",
		"edad":      20,
	})
	if err != nil {
		t.Fatalf("Error insertando datos iniciales: %v", err)
	}

	// Inyectar repo y read-client mock
	services.SetPersonaRepository(&testMongoRepo{col: col})
	services.SetReadClient(readClientMock{})

	// Body del update
	updateData := models.PersonaUpdate{
		Nombre:   "Pedro",
		Apellido: "Ramirez",
		Edad:     30,
	}
	body, _ := json.Marshal(updateData)

	req := httptest.NewRequest("PUT", "/personas/123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/personas/{documento}", ActualizarPersonaHandler).Methods("PUT")

	router.ServeHTTP(rr, req)

	// Validar respuesta HTTP
	if rr.Code != http.StatusOK {
		t.Fatalf("Código esperado 200, recibido %d", rr.Code)
	}

	// Validar actualización en Mongo
	var result models.Persona
	err = col.FindOne(context.Background(), bson.M{"documento": "123"}).Decode(&result)
	if err != nil {
		t.Fatalf("No se pudo obtener persona actualizada: %v", err)
	}

	if result.Nombre != "Pedro" || result.Apellido != "Ramirez" || result.Edad != 30 {
		t.Fatalf("La actualización no se reflejó correctamente en Mongo. Resultado: %+v", result)
	}
}

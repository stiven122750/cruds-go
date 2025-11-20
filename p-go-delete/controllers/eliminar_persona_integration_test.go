//go:build integration

package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stiven122750/cruds-go/p-go-delete/config"
	"github.com/stiven122750/cruds-go/p-go-delete/models"
	"github.com/stiven122750/cruds-go/p-go-delete/repositories"
	"github.com/stiven122750/cruds-go/p-go-delete/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestEliminarPersona_Integration(t *testing.T) {
	ctx := context.Background()

	// =====================================================
	// Levantar MongoDB real en Testcontainers
	// =====================================================
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp").WithStartupTimeout(30 * time.Second),
	}

	mongoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Error iniciando MongoDB en Testcontainers: %v", err)
	}
	defer mongoContainer.Terminate(ctx)

	host, _ := mongoContainer.Host(ctx)
	mappedPort, _ := mongoContainer.MappedPort(ctx, "27017/tcp")

	mongoURI := "mongodb://" + host + ":" + mappedPort.Port()

	// =====================================================
	// Conectar cliente MongoDB real
	// =====================================================
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("Error creando cliente MongoDB: %v", err)
	}
	err = client.Connect(ctx)
	if err != nil {
		t.Fatalf("Error conectando a MongoDB: %v", err)
	}

	db := client.Database("testdb")
	collection := db.Collection("personas")

	// Inyección en el repo real
	config.Collection = collection
	repositories.SetCollection(collection)
	services.SetPersonaRepository(repositories.RealPersonaRepository{})

	// =====================================================
	// Insertar persona real para poder eliminarla
	// =====================================================
	_, err = collection.InsertOne(ctx, bson.M{
		"documento": "555",
		"nombre":    "Persona Test",
	})
	if err != nil {
		t.Fatalf("Error insertando persona previa: %v", err)
	}

	// =====================================================
	// Ejecutar DELETE real contra el controlador
	// =====================================================
	reqHTTP := httptest.NewRequest("DELETE", "/eliminar-persona/555", nil)
	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/eliminar-persona/{documento}", EliminarPersona).Methods("DELETE")
	router.ServeHTTP(w, reqHTTP)

	// =====================================================
	// Validar respuesta HTTP
	// =====================================================
	if w.Code != http.StatusOK {
		t.Fatalf("Se esperaba 200 OK, se obtuvo %d", w.Code)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["mensaje"] != "Persona eliminada exitosamente" {
		t.Fatalf("Mensaje incorrecto: %#v", resp)
	}

	// =====================================================
	// Verificar que se eliminó realmente en MongoDB
	// =====================================================
	var found models.Persona
	err = collection.FindOne(ctx, bson.M{"documento": "555"}).Decode(&found)

	if err == nil {
		t.Fatalf("La persona fue encontrada después de eliminarla — no se eliminó")
	}
}

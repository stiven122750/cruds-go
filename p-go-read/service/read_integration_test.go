package services

import (
	"context"
	"p-go-read/models"
	"p-go-read/repositories"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) *mongo.Collection {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("No se pudo conectar a MongoDB: %v", err)
	}

	db := client.Database("test_read_db_services")
	col := db.Collection("personas")

	// Limpiar antes del test
	col.DeleteMany(context.Background(), bson.M{})

	// Inyectar colección en el repositorio
	repositories.SetCollection(col)

	// Inyectar repositorio real en el service
	SetPersonaRepository(repositories.RealPersonaRepository{})

	return col
}

// -------------------------------------
// TEST INTEGRACIÓN: ObtenerPersona()
// -------------------------------------

func TestIntegration_ObtenerPersona(t *testing.T) {
	col := setupTestDB(t)

	persona := models.Persona{
		Documento: "777",
		Nombre:    "Integration Test",
	}

	_, err := col.InsertOne(context.Background(), persona)
	if err != nil {
		t.Fatalf("No se pudo insertar: %v", err)
	}

	res, err := ObtenerPersona("777")
	if err != nil {
		t.Fatalf("Error consultando: %v", err)
	}

	if res.Documento != "777" {
		t.Fatalf("Documento incorrecto: %s", res.Documento)
	}
}

// -------------------------------------
// TEST INTEGRACIÓN: Listado
// -------------------------------------

func TestIntegration_ObtenerTodas(t *testing.T) {
	col := setupTestDB(t)

	personas := []interface{}{
		models.Persona{Documento: "1", Nombre: "A"},
		models.Persona{Documento: "2", Nombre: "B"},
	}

	_, err := col.InsertMany(context.Background(), personas)
	if err != nil {
		t.Fatalf("Error insertando datos: %v", err)
	}

	list, err := ObtenerTodasLasPersonas()
	if err != nil {
		t.Fatalf("Error consultando personas: %v", err)
	}

	if len(list) != 2 {
		t.Fatalf("Esperaba 2 personas, recibí %d", len(list))
	}
}

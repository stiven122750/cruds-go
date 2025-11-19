package services

import (
	"context"
	"fmt"
	"p-go-read/models"
	"p-go-read/repositories"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ========================
// Setup con Testcontainers
// ========================
func setupTestDB(t *testing.T) *mongo.Collection {

	ctx := context.Background()

	// Crear contenedor MongoDB
	req := testcontainers.ContainerRequest{
		Image:        "mongo:7",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp"),
	}

	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("No se pudo iniciar MongoDB en Testcontainers: %v", err)
	}

	// Obtener puerto asignado por Docker
	endpoint, err := mongoC.Endpoint(ctx, "")
	if err != nil {
		t.Fatalf("Error obteniendo endpoint MongoDB: %v", err)
	}

	// Conexión real
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		fmt.Sprintf("mongodb://%s", endpoint),
	))
	if err != nil {
		t.Fatalf("No se pudo conectar al contenedor MongoDB: %v", err)
	}

	// Tiempo de espera para que Mongo esté listo
	time.Sleep(1 * time.Second)

	db := client.Database("test_read_db_services")
	col := db.Collection("personas")

	// Limpieza
	col.DeleteMany(ctx, bson.M{})

	// Inyectar colección en repositorio real
	repositories.SetCollection(col)
	SetPersonaRepository(repositories.RealPersonaRepository{})

	return col
}

// ========================
// TEST INTEGRACIÓN: ObtenerPersona
// ========================
func TestIntegration_ObtenerPersona(t *testing.T) {
	col := setupTestDB(t)

	ctx := context.Background()

	persona := models.Persona{
		Documento: "777",
		Nombre:    "Integration Test",
	}

	_, err := col.InsertOne(ctx, persona)
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

// ========================
// TEST INTEGRACIÓN: ObtenerTodas
// ========================
func TestIntegration_ObtenerTodas(t *testing.T) {

	col := setupTestDB(t)

	ctx := context.Background()

	personas := []interface{}{
		models.Persona{Documento: "1", Nombre: "A"},
		models.Persona{Documento: "2", Nombre: "B"},
	}

	_, err := col.InsertMany(ctx, personas)
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

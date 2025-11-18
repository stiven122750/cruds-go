//go:build integration
// +build integration

package services

import (
	"context"
	"testing"
	"time"

	"github.com/stiven122750/cruds-go/p-go-delete/config"
	"github.com/stiven122750/cruds-go/p-go-delete/models"
	"github.com/stiven122750/cruds-go/p-go-delete/repositories"
	"github.com/stiven122750/cruds-go/p-go-delete/services"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
)

func TestEliminarPersonaIntegration(t *testing.T) {
	ctx := context.Background()

	// 1️⃣ Crear contenedor de MongoDB
	req := testcontainers.ContainerRequest{
		Image:        "mongo:6.0",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp").WithStartupTimeout(20 * time.Second),
	}
	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	defer mongoC.Terminate(ctx)

	// 2️⃣ Obtener URI del contenedor
	endpoint, err := mongoC.Endpoint(ctx, "")
	assert.NoError(t, err)

	t.Setenv("MONGO_URI", "mongodb://"+endpoint)
	t.Setenv("MONGO_DB", "testdb")
	t.Setenv("COLLECTION_NAME", "personas_test_delete")

	// 3️⃣ Conectar a Mongo
	err = config.ConectarMongo()
	assert.NoError(t, err)
	defer config.CerrarMongo()

	repositories.SetCollection(config.Collection)

	// 4️⃣ Limpiar colección
	_, err = config.Collection.DeleteMany(context.Background(), bson.M{})
	assert.NoError(t, err)

	// 5️Inyectar repositorio real
	services.SetPersonaRepository(repositories.RealPersonaRepository{})

	persona := models.Persona{
		Documento: "99999",
		Nombre:    "Eliminar",
		Apellido:  "Test",
		Edad:      30,
		Correo:    "delete@test.com",
		Telefono:  "3000000000",
		Direccion: "Sin dirección",
	}

	err = services.CrearPersona(persona)
	assert.NoError(t, err)

	// 7️⃣ Ejecutar eliminación
	err = services.EliminarPersona(persona.Documento)
	assert.NoError(t, err)

	// 8️⃣ Verificar que ya no exista
	var resultado models.Persona
	err = config.Collection.FindOne(ctx, bson.M{"documento": persona.Documento}).Decode(&resultado)
	assert.Error(t, err) // Debe fallar porque no existe
}

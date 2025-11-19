//go:build integration
// +build integration

package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/stiven122750/cruds-go/p-go-delete/config"
	"github.com/stiven122750/cruds-go/p-go-delete/models"
	"github.com/stiven122750/cruds-go/p-go-delete/repositories"
	"github.com/stiven122750/cruds-go/p-go-delete/services"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"go.mongodb.org/mongo-driver/bson"
)

func TestEliminarPersonaIntegration(t *testing.T) {
	ctx := context.Background()

	// 1️⃣ Crear contenedor MongoDB
	req := testcontainers.ContainerRequest{
		Image:        "mongo:6.0",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor: wait.
			ForListeningPort("27017/tcp").
			WithStartupTimeout(40 * time.Second), // GH Actions es lento
	}

	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	defer mongoC.Terminate(ctx)

	// 2️⃣ Obtener endpoint real del contenedor
	endpoint, err := mongoC.Endpoint(ctx, "27017/tcp")
	assert.NoError(t, err)

	// 3️⃣ Inyectar variables de entorno para config.Mongo
	t.Setenv("MONGO_URI", "mongodb://"+endpoint)
	t.Setenv("MONGO_DB", "testdb_delete")
	t.Setenv("COLLECTION_NAME", "personas_test_delete")

	// 4️⃣ Conectar Mongo
	err = config.ConectarMongo()
	assert.NoError(t, err)
	defer config.CerrarMongo()

	// 5️⃣ Enviar colección al repositorio real
	repositories.SetCollection(config.Collection)

	// 6️⃣ Limpiar colección antes del test
	_, err = config.Collection.DeleteMany(ctx, bson.M{})
	assert.NoError(t, err)

	// 7️⃣ Inyectar repositorio real en el service DELETE
	services.SetPersonaRepository(repositories.RealPersonaRepository{})

	// 8️⃣ Insertar persona a eliminar
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

	// 9️⃣ Eliminar persona
	err = services.EliminarPersona(persona.Documento)
	assert.NoError(t, err)

	// 🔟 Validar que ya no exista
	var resultado models.Persona
	err = config.Collection.FindOne(ctx, bson.M{"documento": persona.Documento}).Decode(&resultado)

	assert.Error(t, err) // Debe fallar: ya no existe
}

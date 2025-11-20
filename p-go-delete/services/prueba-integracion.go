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

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"go.mongodb.org/mongo-driver/bson"
)

func TestEliminarPersonaIntegration(t *testing.T) {
	ctx := context.Background()

	// 1️⃣ Crear contenedor MongoDB real
	req := testcontainers.ContainerRequest{
		Image:        "mongo:6.0",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor: wait.ForListeningPort("27017/tcp").
			WithStartupTimeout(40 * time.Second),
	}

	mongoC, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	assert.NoError(t, err)
	defer mongoC.Terminate(ctx)

	// 2️⃣ Obtener endpoint real (host:port)
	endpoint, err := mongoC.Endpoint(ctx, "27017/tcp")
	assert.NoError(t, err)

	// 3️⃣ Setear variables de entorno para config.ConectarMongo()
	t.Setenv("MONGO_URI", "mongodb://"+endpoint)
	t.Setenv("MONGO_DB", "testdb_delete")
	t.Setenv("COLLECTION_NAME", "personas_test_delete")

	// 4️⃣ Conectar a Mongo real
	err = config.ConectarMongo()
	assert.NoError(t, err)

	// 5️⃣ Inyectar colección real en el repositorio
	repositories.SetCollection(config.Collection)

	// 6️⃣ Borrar colección antes del test
	_, err = config.Collection.DeleteMany(ctx, bson.M{})
	assert.NoError(t, err)

	// 7️⃣ Inyectar repositorio real en los servicios
	SetPersonaRepository(repositories.RealPersonaRepository{})

	// 8️⃣ Insertar persona antes de eliminar (NO existe CreatePersona)
	persona := models.Persona{
		Documento: "99999",
		Nombre:    "Eliminar",
		Apellido:  "Test",
	}

	_, err = config.Collection.InsertOne(ctx, persona)
	assert.NoError(t, err)

	// 9️⃣ Ejecutar eliminación real
	err = EliminarPersona(persona.Documento)
	assert.NoError(t, err)

	// 🔟 Validar que ya no exista en Mongo
	var result models.Persona
	err = config.Collection.FindOne(ctx, bson.M{"documento": persona.Documento}).Decode(&result)

	assert.Error(t, err) // OK: ya debe estar eliminado
}

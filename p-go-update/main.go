package main

import (
	"fmt"
	"log"
	"net/http"
	"p-go-update/config"
	"p-go-update/controllers"
	"p-go-update/repositories"
	services "p-go-update/service"

	"github.com/gorilla/mux"
)

func main() {

	// Conectar a MongoDB
	err := config.ConectarMongo()
	if err != nil {
		log.Fatal("Error conectando a MongoDB:", err)
	}

	// Inyectar repositorio real
	services.SetPersonaRepository(repositories.RealPersonaRepository{})

	// Inyectar colección Mongo
	repositories.SetCollection(config.Collection)

	// ❗ INYECTAR EL CLIENTE READ REAL (FALTABA ESTO)
	services.SetReadClient(services.ReadServiceClient{})

	// Router
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hola, este es el microservicio de actualización de personas")
	})

	router.HandleFunc("/personas/{documento}", controllers.ActualizarPersonaHandler).Methods("PUT")

	puerto := ":5001"
	fmt.Printf("API UPDATE escuchando en http://localhost%s\n", puerto)
	log.Fatal(http.ListenAndServe(puerto, router))
}

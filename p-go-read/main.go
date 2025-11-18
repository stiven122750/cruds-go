package main

import (
	"fmt"
	"log"
	"net/http"
	services "p-go-read/service"

	"p-go-read/config"
	"p-go-read/controllers"
	"p-go-read/repositories"

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

	// Inyectar colección de MongoDB
	repositories.SetCollection(config.Collection)

	// Router
	router := mux.NewRouter()

	// Mensaje de bienvenida del microservicio
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hola, este es el microservicio de lectura de personas")
	})

	// Rutas READ
	router.HandleFunc("/personas/{documento}", controllers.ObtenerPersonaHandler).Methods("GET")
	router.HandleFunc("/personas", controllers.ObtenerTodasLasPersonasHandler).Methods("GET")

	// Puerto del microservicio READ
	puerto := ":5000"
	fmt.Printf("API READ escuchando en http://localhost%s\n", puerto)
	log.Fatal(http.ListenAndServe(puerto, router))
}

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

	// Conectar a MongoDB (si lo necesitas para actualizar)
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
		fmt.Fprintln(w, "Hola, este es el microservicio de actualización de personas")
	})

	// Rutas UPDATE
	router.HandleFunc("/personas/{documento}", controllers.ActualizarPersonaHandler).Methods("PUT")

	// Puerto del microservicio UPDATE
	puerto := ":5001" // distinto del READ
	fmt.Printf("API UPDATE escuchando en http://localhost%s\n", puerto)
	log.Fatal(http.ListenAndServe(puerto, router))
}

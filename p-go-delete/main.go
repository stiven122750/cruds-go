package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/stiven122750/cruds-go/p-go-delete/config"
	"github.com/stiven122750/cruds-go/p-go-delete/controllers"
	"github.com/stiven122750/cruds-go/p-go-delete/repositories"
	"github.com/stiven122750/cruds-go/p-go-delete/services"

	"github.com/gorilla/mux"
)

func main() {
	// Conectar a MongoDB
	err := config.ConectarMongo()
	if err != nil {
		log.Fatal("❌ Error conectando a MongoDB:", err)
	}

	// Inyectar repositorio real
	services.SetPersonaRepository(repositories.RealPersonaRepository{})

	// Inyectar colección de MongoDB
	repositories.SetCollection(config.Collection)

	// Router
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hola! API para eliminar personas en MongoDB")
	})

	// Ruta DELETE
	router.HandleFunc("/eliminar-persona/{documento}", controllers.EliminarPersona).Methods("DELETE")

	// Puerto
	puerto := ":8080"
	fmt.Printf("🚀 API DELETE lista en http://localhost%s\n", puerto)
	log.Fatal(http.ListenAndServe(puerto, router))
}

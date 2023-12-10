package main

import (
	"log"
	"main-api/db"
	"net/http"
)

import (
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"main-api/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	r := chi.NewRouter()

	db.ConnectDatabase()

	r.Mount("/api", routes.ApiRouter())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("404"))
	})

	log.Println("Running on localhost:8080")
	http.ListenAndServe(":8080", r)

}

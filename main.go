package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Big Error , ENV NOT LOADED %w", err)
		return
	}

	db, e := sql.Open("sqlite3", "./"+os.Getenv("EVENT_DB"))
	Handle_error(e)
	defer db.Close()
	CheckDB(db)

	r := chi.NewRouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}).Handler)
	r.Route("/api", func(r chi.Router) {
		r.Get("/events", GetFromDB(db))
		r.Post("/events", InsertDB(db))
	})
	fmt.Println("Listening on port http://localhost:8080")
	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var data = []User{
	{ID: 1, Name: "Ashwin"},
	{ID: 2, Name: "ME"},
}

func main() {
	r := chi.NewRouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}).Handler)
	r.Route("/api", func(r chi.Router) {
		r.Get("/users", getUser)
		r.Get("/projects",onPro)
	})

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", r)
}

func onPro(w http.ResponseWriter, r *http.Request){
	fmt.Println("GET CALLED")
	json.NewEncoder(w).Encode(data)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET CALLED")
	json.NewEncoder(w).Encode(data)
}

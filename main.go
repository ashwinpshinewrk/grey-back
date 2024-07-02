package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/mattn/go-sqlite3"
)

func handle_error(e error) {
	if e != nil {
		fmt.Println(e)
		return
	}
}

type Event struct {
	ID          int    `json:"id"`
	Image       []byte `json:"event_image"`
	Title       string `json:"event_title"`
	Description string `json:"event_description"`
}

func (e Event) new(id int, image []byte, title string, description string) Event {
	var event Event
	event.ID = id
	event.Image = image
	event.Title = title
	event.Description = description
	return event
}

var db *sql.DB

func checkDB(db *sql.DB) {
	var e error
	stmt, e := db.Prepare("create table if not exists event (id integer PRIMARY KEY,event_image blob,event_title text,event_description text);")
	handle_error(e)
	stmt.Exec()
}

// not to use
func getAllEvent(db *sql.DB) ([]Event, error) {
	rows, err := db.Query("select id, event_image, event_title, event_description from event;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err = rows.Scan(&event.ID, &event.Image, &event.Title, &event.Description)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func getFromDB(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		events, err := getAllEvent(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(events)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

func insertDB(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event Event
		var e error

		e = json.NewDecoder(r.Body).Decode(&event)
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
		}
		stmt, e := db.Prepare("insert into event (id, event_image, event_title, event_description) values (?,?,?,?); ")
		if e != nil {
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}
		stmt.Exec(event.ID, event.Image, event.Title, event.Description)
		w.WriteHeader(http.StatusCreated)
	}
}

func main() {
	db, e := sql.Open("sqlite3", "./events.db")
	handle_error(e)
	defer db.Close()
	checkDB(db)

	r := chi.NewRouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}).Handler)
	r.Route("/api", func(r chi.Router) {
		r.Get("/events", getFromDB(db))
		r.Post("/events", insertDB(db))
	})
	// r.Route("/api", func(r chi.Router) {
	// 	r.Get("/users", getUser)
	// 	r.Get("/projects",onPro)
	// })

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", r)
}

// func onPro(w http.ResponseWriter, r *http.Request){
// 	fmt.Println("GET CALLED")
// 	json.NewEncoder(w).Encode(data)
// }

// func getUser(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("GET CALLED")
// 	json.NewEncoder(w).Encode(data)
// }

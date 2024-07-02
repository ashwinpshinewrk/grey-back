package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Creates a table if not exists
func CheckDB(db *sql.DB) {
	var e error
	stmt, e := db.Prepare("create table if not exists event (id integer PRIMARY KEY,event_image blob,event_title text,event_description text);")
	Handle_error(e)
	stmt.Exec()
}

// Extracts all data from the table
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

// A helper for GET METHOD for data extraction
func GetFromDB(db *sql.DB) http.HandlerFunc {
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

// Inserts data into the DB
func InsertDB(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event Event
		var e error

		e = json.NewDecoder(r.Body).Decode(&event)
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
		}
		fmt.Println(event.Pass)
		if event.Pass != os.Getenv("PASS") {
			http.Error(w, "NOT VALID PASS", http.StatusBadRequest)
			return
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

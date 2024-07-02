package main

// Declaration of Event JSON Object
type Event struct {
	ID          int    `json:"id"`
	Image       []byte `json:"event_image"`
	Title       string `json:"event_title"`
	Description string `json:"event_description"`
	Pass        string `json:"pass"`
}

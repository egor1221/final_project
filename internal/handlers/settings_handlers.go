package handlers

import (
	"os"
	"time"
)

var now time.Time = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())

var webDir string = "web"

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type password struct {
	Password string `json:"password"`
}

var pass string = os.Getenv("TODO_PASSWORD")

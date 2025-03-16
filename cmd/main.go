package main

import (
	"final_project/internal/database"
	"final_project/internal/handlers"
	"log"
	"net/http"
	"os"
)

var todoPort string = os.Getenv("TODO_PORT")

func main() {

	r := handlers.Router()

	database.CheckDb()

	if len(todoPort) == 0 {
		todoPort = ":7540"
	}

	err := http.ListenAndServe(todoPort, r)

	if err != nil {
		log.Fatalf(err.Error())
	}
}

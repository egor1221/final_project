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

	db, err := database.OpenDB()

	if err != nil {
		log.Fatalf(err.Error())
	}

	defer db.Close()

	r := handlers.Router(db)

	if len(todoPort) == 0 {
		todoPort = ":7540"
	}

	log.Println("Порт запуска: " + todoPort)

	err = http.ListenAndServe(todoPort, r)

	if err != nil {
		log.Fatalf(err.Error())
	}
}

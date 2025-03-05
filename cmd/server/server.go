package server

import (
	"net/http"
	"os"
	"todo/database"
)

var webDir string = "../web"

var todoPort string = os.Getenv("TODO_PORT")

func StartServer() error {

	db, err := database.OpenDB()

	if err != nil {
		return err
	}
	defer db.Close()

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	err = http.ListenAndServe(todoPort, nil)

	return err

}

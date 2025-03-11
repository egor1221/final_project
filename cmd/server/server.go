package server

import (
	"net/http"
	"os"
	"todo/database"

	"github.com/go-chi/chi"
)

var webDir string = "../web"

var todoPort string = os.Getenv("TODO_PORT")

func StartServer() error {

	db, err := database.OpenDB()

	if err != nil {
		return err
	}
	defer db.Close()

	r := chi.NewRouter()

	r.Get("/api/nextdate", getRepeat)
	r.Post("/api/task", postTask)
	r.Get("/api/task", getTask)

	r.Handle("/", http.FileServer(http.Dir(webDir)))
	r.Handle("/css/style.css", http.FileServer(http.Dir(webDir)))
	r.Handle("/css/theme.css", http.FileServer(http.Dir(webDir)))
	r.Handle("/js/axios.min.js", http.FileServer(http.Dir(webDir)))
	r.Handle("/js/scripts.min.js", http.FileServer(http.Dir(webDir)))
	r.Handle("/favicon.ico", http.FileServer(http.Dir(webDir)))

	err = http.ListenAndServe(todoPort, r)

	return err
}

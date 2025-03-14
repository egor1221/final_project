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

	if len(todoPort) == 0 {
		todoPort = ":7540"
	}

	r := chi.NewRouter()

	r.Get("/api/nextdate", getRepeat)

	r.Get("/api/tasks", getTasks)

	r.Get("/api/task", getTaskById)
	r.Post("/api/task", postTask)
	r.Put("/api/task", putTask)

	r.Post("/api/task/done", postCheck)
	r.Delete("/api/task", deleteTask)

	r.Handle("/", http.FileServer(http.Dir(webDir)))
	r.Handle("/css/style.css", http.FileServer(http.Dir(webDir)))
	r.Handle("/css/theme.css", http.FileServer(http.Dir(webDir)))
	r.Handle("/js/axios.min.js", http.FileServer(http.Dir(webDir)))
	r.Handle("/js/scripts.min.js", http.FileServer(http.Dir(webDir)))
	r.Handle("/favicon.ico", http.FileServer(http.Dir(webDir)))
	r.Handle("/login.html", http.FileServer(http.Dir(webDir)))
	r.Handle("/index.html", http.FileServer(http.Dir(webDir)))

	err = http.ListenAndServe(todoPort, r)

	return err
}

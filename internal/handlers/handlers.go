package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Router() *chi.Mux {
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

	return r
}

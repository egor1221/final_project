package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Router() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/api/nextdate", getRepeat)

	r.Get("/api/tasks", authMiddleware(getTasks))

	r.Get("/api/task", authMiddleware(getTaskById))
	r.Post("/api/task", authMiddleware(postTask))
	r.Put("/api/task", authMiddleware(putTask))

	r.Post("/api/task/done", authMiddleware(postCheck))
	r.Delete("/api/task", authMiddleware(deleteTask))

	r.Post("/api/signin", postPassword)

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

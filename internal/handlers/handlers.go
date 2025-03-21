package handlers

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
)

func Router(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/api/nextdate", getRepeat)

	r.Get("/api/tasks", authMiddleware(getTasks(db)))

	r.Get("/api/task", authMiddleware(getTaskById(db)))
	r.Post("/api/task", authMiddleware(postTask(db)))
	r.Put("/api/task", authMiddleware(putTask(db)))

	r.Post("/api/task/done", authMiddleware(postCheck(db)))
	r.Delete("/api/task", authMiddleware(deleteTask(db)))

	r.Post("/api/signin", postPassword)

	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(webDir))))

	return r
}

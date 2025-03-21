package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"final_project/internal/database"
	"final_project/internal/repeattask"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func getRepeat(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	timeNow, err := time.Parse("20060102", now)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	nextDate, err := repeattask.NextDate(timeNow, date, repeat)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(nextDate))
}

func postTask(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)

		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
			http.Error(w, `{"error": "ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		if len(task.Title) == 0 {
			http.Error(w, `{"error": "не указан заголовок задачи"}`, http.StatusBadRequest)
			return
		}

		if task.Date == "" {
			task.Date = now.Format("20060102")
		}

		t, err := time.Parse("20060102", task.Date)

		if err != nil {
			http.Error(w, `{"error": "ошибка выполнения time.Parse, дата должна быть в формате YYYYMMDD"}`, http.StatusBadRequest)
			return
		}

		if t.Format("20060102") != task.Date {
			http.Error(w, `{"error": "дата представлена в формате, отличном от 20060102"}`, http.StatusBadRequest)
			return
		}

		if t.Before(now) && task.Date != now.Format("20060102") {
			if task.Repeat != "" {

				repeatTask, err := repeattask.NextDate(now, task.Date, task.Repeat)

				if err != nil {
					http.Error(w, `{"error": "правило повторения указано в неправильном формате"}`, http.StatusBadRequest)
					return
				}
				task.Date = repeatTask
			} else {
				task.Date = now.Format("20060102")
			}
		}

		id, err := database.AddTask(db, task.Date, task.Title, task.Comment, task.Repeat)

		if err != nil {
			http.Error(w, `{"error": "ошибка при добавлении задачи"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(map[string]int64{"id": id})
	}
}

func getTasks(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")

		if date, err := time.Parse("02.01.2006", search); err == nil {
			search = date.Format("20060102")
		}

		tasks := map[string][]Task{
			"tasks": {},
		}

		rows, err := database.SelectTasks(db, search)

		if err != nil {
			http.Error(w, `{"error": "ошибка при получении данных"}`, http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			task := Task{}

			err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
				return
			}

			tasks["tasks"] = append(tasks["tasks"], task)
		}

		err = rows.Err()

		if err != nil {
			http.Error(w, `{"error": "ошибка при получении данных"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(tasks)
	}
}

func getTaskById(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var task Task

		id := r.URL.Query().Get("id")

		if len(id) == 0 {
			http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
			return
		}

		row, err := database.SelectTaskById(db, id)

		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		if row == nil {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusInternalServerError)
			return
		}

		err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(task)
	}
}

func putTask(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)

		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
			http.Error(w, `{"error": "ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		if len(task.Title) == 0 {
			http.Error(w, `{"error": "не указан заголовок задачи"}`, http.StatusBadRequest)
			return
		}

		if task.Date == "" {
			task.Date = now.Format("20060102")
		}

		t, err := time.Parse("20060102", task.Date)

		if err != nil || t.Format("20060102") != task.Date {
			http.Error(w, `{"error": "дата представлена в формате, отличном от 20060102"}`, http.StatusBadRequest)
			return
		}

		if t.Before(now) {
			http.Error(w, `{"error": "дата не должна быть позже текущей"}`, http.StatusBadRequest)
			return
		}

		res, err := database.UpdateTask(db, task.ID, task.Date, task.Title, task.Comment, task.Repeat)

		if err != nil || res == 0 {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(map[string]string{})
	}
}

func postCheck(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		id := r.URL.Query().Get("id")

		row, err := database.SelectTaskById(db, id)

		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		if row == nil {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusInternalServerError)
			return
		}

		err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		if len(task.Repeat) == 0 {
			err := database.DeleteTask(db, task.ID)

			if err != nil {
				http.Error(w, `{"error": "ошибка при удалении"}`, http.StatusInternalServerError)
				return
			}

		} else {
			repeatTask, err := repeattask.NextDate(now, task.Date, task.Repeat)

			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
				return
			}
			task.Date = repeatTask

			_, err = database.UpdateTask(db, task.ID, task.Date, task.Title, task.Comment, task.Repeat)

			if err != nil {
				http.Error(w, `{"error": "Задача не найдена"}`, http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(map[string]string{})
	}
}

func deleteTask(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		if id == "" {
			http.Error(w, `{"error": "id пустой"}`, http.StatusBadRequest)
			return
		}

		parseId, err := strconv.Atoi(id)

		if err != nil || parseId <= 0 {
			http.Error(w, `{"error": "неверный формат id"}`, http.StatusBadRequest)
			return
		}

		err = database.DeleteTask(db, id)

		if err != nil {
			http.Error(w, `{"error": "ошибка при удалении"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(map[string]string{})
	}
}

func postPassword(w http.ResponseWriter, r *http.Request) {
	var password password
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &password); err != nil {
		http.Error(w, `{"error": "ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	}

	if pass != password.Password {
		http.Error(w, `{"error": "неверный пароль"}`, http.StatusUnauthorized)
		return
	}

	signedToken, err := generateToken()

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{"token": signedToken})
}

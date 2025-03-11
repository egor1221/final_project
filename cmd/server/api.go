package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"todo/database"
	"todo/repeattask"
)

type Task struct {
	ID      int    `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func getRepeat(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	timeNow, err := time.Parse("20060102", now)

	if err != nil {
		fmt.Println(err)
		return
	}

	repeatTask, err := repeattask.NextDate(timeNow, date, repeat)

	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write([]byte(repeatTask))
}

func postTask(w http.ResponseWriter, r *http.Request) {
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
		task.Date = time.Now().Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)

	if err != nil || t.Format("20060102") != task.Date {
		http.Error(w, `{"error": "дата представлена в формате, отличном от 20060102"}`, http.StatusBadRequest)
		return
	}

	if task.Repeat == "" {
		task.Date = time.Now().Format("20060102")
	} else {
		repeatTask, err := repeattask.NextDate(time.Now(), task.Date, task.Repeat)

		if err != nil {
			http.Error(w, `{"error": "правило повторения указано в неправильном формате"}`, http.StatusBadRequest)
			return
		}
		task.Date = repeatTask
	}

	id, err := database.AddTask(task.Date, task.Title, task.Comment, task.Repeat)

	if err != nil {
		fmt.Println(err)
		http.Error(w, `{"error": "ошибка при добавлении задачи"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"id": "%d"}`, id)))
}

func getTask(w http.ResponseWriter, r *http.Request) {
	var tasks []Task

	rows, err := database.SelectTasks()

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		tasks = append(tasks, task)
	}

	result, err := json.Marshal(tasks)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(result))
}

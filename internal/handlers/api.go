package handlers

import (
	"bytes"
	"encoding/json"
	"final_project/internal/database"
	"final_project/internal/repeattask"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Task struct {
	ID      string `json:"id,omitempty"`
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
	if _, err := w.Write([]byte(nextDate)); err != nil {
		fmt.Println(err.Error())
	}
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

	if err != nil {
		http.Error(w, `{"error": "ошибка выполнения time.Parse, дата должна быть в формате YYYYMMDD"}`, http.StatusBadRequest)
		return
	}

	if t.Format("20060102") != task.Date {
		http.Error(w, `{"error": "дата представлена в формате, отличном от 20060102"}`, http.StatusBadRequest)
		return
	}

	if task.Repeat != "" && t.Before(time.Now()) {

		repeatTask, err := repeattask.NextDate(time.Now(), task.Date, task.Repeat)

		if err != nil {
			http.Error(w, `{"error": "правило повторения указано в неправильном формате"}`, http.StatusBadRequest)
			return
		}
		task.Date = repeatTask
	} else if time.Now().Format("20060102") == task.Date {
		task.Date = time.Now().Format("20060102")
	}

	id, err := database.AddTask(task.Date, task.Title, task.Comment, task.Repeat)

	if err != nil {
		http.Error(w, `{"error": "ошибка при добавлении задачи"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write([]byte(fmt.Sprintf(`{"id": "%d"}`, id))); err != nil {
		fmt.Println(err.Error())
	}
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	if date, err := time.Parse("02.01.2006", search); err == nil {
		search = date.Format("20060102")
	}

	tasks := map[string][]Task{
		"tasks": {},
	}

	rows, err := database.SelectTasks(search)

	if err != nil {
		fmt.Println(err)
		http.Error(w, `{"error": "ошибка при получении данных"}`, http.StatusBadRequest)
		return
	}

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		tasks["tasks"] = append(tasks["tasks"], task)
	}

	result, err := json.Marshal(tasks)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(result); err != nil {
		fmt.Println(err.Error())
	}
}

func getTaskById(w http.ResponseWriter, r *http.Request) {
	var task Task

	id := r.URL.Query().Get("id")

	if len(id) == 0 {
		http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	row := database.SelectTaskById(id)

	if row == nil {
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusBadRequest)
		return
	}

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(task)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(result); err != nil {
		fmt.Println(err.Error())
	}

}

func putTask(w http.ResponseWriter, r *http.Request) {
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

	if task.Repeat != "" {
		repeatTask, err := repeattask.NextDate(time.Now(), task.Date, task.Repeat)

		if err != nil {
			http.Error(w, `{"error": "правило повторения указано в неправильном формате"}`, http.StatusBadRequest)
			return
		}
		task.Date = repeatTask
	}

	res, err := database.UpdateTask(task.ID, task.Date, task.Title, task.Comment, task.Repeat)

	if err != nil || res == 0 {
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write([]byte("{}")); err != nil {
		fmt.Println(err.Error())
	}
}

func postCheck(w http.ResponseWriter, r *http.Request) {
	var task Task
	id := r.URL.Query().Get("id")

	row := database.SelectTaskById(id)

	if row == nil {
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusBadRequest)
		return
	}

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	if len(task.Repeat) == 0 {
		err := database.DeleteTask(task.ID)

		if err != nil {
			http.Error(w, `{"error": "ошибка при удалении"}`, http.StatusBadRequest)
			return
		}

	} else {
		repeatTask, err := repeattask.NextDate(time.Now(), task.Date, task.Repeat)

		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		task.Date = repeatTask

		_, err = database.UpdateTask(task.ID, task.Date, task.Title, task.Comment, task.Repeat)

		if err != nil {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write([]byte("{}")); err != nil {
		fmt.Println(err.Error())
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
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

	err = database.DeleteTask(id)

	if err != nil {
		http.Error(w, `{"error": "ошибка при удалении"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write([]byte("{}")); err != nil {
		fmt.Println(err.Error())
	}
}

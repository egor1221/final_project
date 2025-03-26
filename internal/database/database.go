package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func checkDb() {

	if dbFile != "" {
		return
	}

	appPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dbFile = filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	if err != nil {
		err := createTable()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func OpenDB() (*sql.DB, error) {
	checkDb()
	return sql.Open("sqlite", dbFile)
}

func createTable() error {

	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(schemeSql)

	if err != nil {
		return err
	}

	return nil

}

func AddTask(db *sql.DB, date, title, comment, repeat string) (int64, error) {

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", date),
		sql.Named("title", title),
		sql.Named("comment", comment),
		sql.Named("repeat", repeat))

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return id, err
}

func SelectTasks(db *sql.DB, search string) (*sql.Rows, error) {

	var rows *sql.Rows
	var err error

	if search != "" {
		search = "%" + search + "%"
		rows, err = db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :search OR date LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit",
			sql.Named("limit", limit),
			sql.Named("search", search))

		if err != nil {
			return nil, err
		}
	} else {
		rows, err = db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit",
			sql.Named("limit", limit),
			sql.Named("search", search))

		if err != nil {
			return nil, err
		}
	}

	return rows, err
}

func SelectTaskById(db *sql.DB, id string) (*sql.Row, error) {

	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id=:id",
		sql.Named("id", id))

	return row, nil
}

func UpdateTask(db *sql.DB, id, date, title, comment, repeat string) (int64, error) {

	res, err := db.Exec("UPDATE scheduler SET date=:date, title=:title, comment=:comment, repeat=:repeat WHERE id=:id",
		sql.Named("date", date),
		sql.Named("title", title),
		sql.Named("comment", comment),
		sql.Named("repeat", repeat),
		sql.Named("id", id))

	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func DeleteTask(db *sql.DB, id string) error {

	_, err := db.Exec("DELETE FROM scheduler WHERE id=:id", sql.Named("id", id))

	if err != nil {
		return err
	}

	return nil
}

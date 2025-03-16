package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var dbFile string = os.Getenv("TODO_DBFILE")
var limit int = 50

func CheckDb() {

	if dbFile != "" {
		fmt.Println(123)
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

func openDB() (*sql.DB, error) {

	db, err := sql.Open("sqlite", dbFile)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTable() error {
	schemeSql := `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(128) NOT NULL DEFAULT "",
    comment VARCHAR(128) NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
    );
    CREATE INDEX scheduler_date ON scheduler (date)`

	db, err := sql.Open("sqlite", "../scheduler.db")
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

func AddTask(date, title, comment, repeat string) (int64, error) {
	db, err := openDB()

	if err != nil {
		return 0, err
	}

	defer db.Close()

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", date),
		sql.Named("title", title),
		sql.Named("comment", comment),
		sql.Named("repeat", repeat))

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return id, err
}

func SelectTasks(search string) (*sql.Rows, error) {
	db, err := openDB()

	if err != nil {
		return nil, err
	}

	defer db.Close()

	var rows *sql.Rows

	if search != "" {
		search += "%"
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

	return rows, nil
}

func SelectTaskById(id string) *sql.Row {
	db, err := openDB()

	if err != nil {
		log.Fatalf(err.Error())
		return nil
	}
	defer db.Close()

	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id=:id",
		sql.Named("id", id))

	return row
}

func UpdateTask(id, date, title, comment, repeat string) (int64, error) {
	db, err := openDB()

	if err != nil {
		return 0, err
	}
	defer db.Close()

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

func DeleteTask(id string) error {
	db, err := openDB()

	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM scheduler WHERE id=:id", sql.Named("id", id))

	if err != nil {
		return err
	}

	return nil
}

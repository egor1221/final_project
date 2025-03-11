package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var install bool

func checkDb() {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	if err != nil {
		install = true
	}
}

func OpenDB() (*sql.DB, error) {
	checkDb()

	if install {
		err := createTable()
		if err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite", "scheduler.db")

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

	file, err := os.Create("scheduler.db")

	if err != nil {
		return err
	}
	defer file.Close()

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(schemeSql)

	if err != nil {
		return err
	}

	install = false

	return nil

}

func AddTask(date, title, comment, repeat string) (int64, error) {
	db, err := OpenDB()

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
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return id, err
}

func SelectTasks() (*sql.Rows, error) {
	db, err := OpenDB()

	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date")

	if err != nil {
		return nil, err
	}

	return rows, nil
}

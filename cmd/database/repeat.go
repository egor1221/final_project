package database

import (
	"fmt"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	db, err := OpenDB()

	if err != nil {
		return "", err
	}

	defer db.Close()

	if len(repeat) == 0 {
		_, err := db.Exec("DELETE FROM scheduler WHERE date = :date", date)

		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("в колонке repeat — пустая строка")
	}
}

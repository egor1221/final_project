package repeattask

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	// "todo/database"
)

var letter string
var number int

func NextDate(now time.Time, date string, repeat string) (string, error) {
	// db, err := database.OpenDB()

	// if err != nil {
	// 	return "", err
	// }

	// defer db.Close()

	// if len(repeat) == 0 {
	// 	_, err := db.Exec("DELETE FROM scheduler WHERE date = :date", date)

	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	return "", fmt.Errorf("в колонке repeat — пустая строка")
	// }

	parseDate, err := time.Parse("20060102", date)

	if err != nil {
		return "", fmt.Errorf(`время в переменной date не может быть преобразовано в корректную дату — ошибка выполнения time.Parse("20060102", d)`)
	}

	letter, number = splitRepeat(repeat)

	var transformDate string

	switch letter {
	case "y":
		transformDate = addYear(now, parseDate)
	case "d":
		if number == 0 && number > 400 {
			return "", fmt.Errorf("указан неверный формат repeat")
		}
		transformDate = addDay(now, parseDate, number)
	default:
		return "", fmt.Errorf("указан неверный формат repeat")
	}

	return transformDate, nil
}

func splitRepeat(repeat string) (string, int) {
	arr := strings.Split(repeat, " ")

	if len(arr) < 2 {
		return arr[0], 0
	}

	num, err := strconv.Atoi(arr[1])

	if err != nil {
		fmt.Println(err)
	}

	return arr[0], num
}

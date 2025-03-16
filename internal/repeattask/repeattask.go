package repeattask

import (
	"fmt"
	"time"
)

var letter string

func NextDate(now time.Time, date string, repeat string) (string, error) {

	parseDate, err := time.Parse("20060102", date)

	if err != nil {
		return "", fmt.Errorf(`время в переменной date не может быть преобразовано в корректную дату — ошибка выполнения time.Parse("20060102", d)`)
	}

	if repeat == "" {
		return "", fmt.Errorf(`в колонке repeat — пустая строка`)
	}

	letter = string(repeat[0])

	var nextDate string

	switch letter {
	case "y":
		nextDate = addYear(now, parseDate)
	case "d":
		nextDate, err = addDay(now, parseDate, repeat)
		if err != nil {
			return "", err
		}
	case "w":
		nextDate, err = addWeek(now, parseDate, repeat)
		if err != nil {
			return "", err
		}
	case "m":
		nextDate, err = addMonth(now, parseDate, repeat)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("указан неверный формат repeat")
	}

	return nextDate, nil
}

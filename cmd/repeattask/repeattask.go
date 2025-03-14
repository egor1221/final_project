package repeattask

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var letter string
var number int

func NextDate(now time.Time, date string, repeat string) (string, error) {

	parseDate, err := time.Parse("20060102", date)

	if err != nil {
		return "", fmt.Errorf(`время в переменной date не может быть преобразовано в корректную дату — ошибка выполнения time.Parse("20060102", d)`)
	}

	if repeat == "" {
		return "", fmt.Errorf(`в колонке repeat — пустая строка`)
	}

	letter, number = splitRepeat(repeat)

	var nextDate string

	switch letter {
	case "y":
		nextDate = addYear(now, parseDate)
	case "d":
		nextDate, err = addDay(now, parseDate, number)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("указан неверный формат repeat")
	}

	return nextDate, nil
}

func splitRepeat(repeat string) (string, int) {
	arr := strings.Fields(repeat)

	if len(arr) < 2 {
		return arr[0], 0
	}

	num, err := strconv.Atoi(arr[1])

	if err != nil {
		fmt.Println(err)
	}

	return arr[0], num
}

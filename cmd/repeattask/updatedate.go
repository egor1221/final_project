package repeattask

import (
	"fmt"
	"time"
)

func addYear(now, parseDate time.Time) string {

	resDate := parseDate.AddDate(1, 0, 0)

	for resDate.Before(now) {
		resDate = resDate.AddDate(1, 0, 0)
	}

	return resDate.Format("20060102")
}

func addDay(now, parseDate time.Time, number int) (string, error) {

	if number == 0 || number > 400 {
		return "", fmt.Errorf("указан неверный формат repeat")
	}

	resDate := parseDate.AddDate(0, 0, number)
	fmt.Println(resDate)

	for resDate.Before(now) {
		resDate = resDate.AddDate(0, 0, number)
	}

	return resDate.Format("20060102"), nil
}

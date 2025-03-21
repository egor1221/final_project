package repeattask

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func addYear(now, parseDate time.Time) string {

	resDate := parseDate.AddDate(1, 0, 0)

	for resDate.Before(now) {
		resDate = resDate.AddDate(1, 0, 0)
	}

	return resDate.Format("20060102")
}

func addDay(now, parseDate time.Time, repeat string) (string, error) {

	repeatArr := strings.Fields(repeat)

	if len(repeatArr) < 2 {
		return "", fmt.Errorf("указан неверный формат repeat")
	}

	number, err := strconv.Atoi(repeatArr[1])

	if err != nil || number == 0 || number > 400 {
		return "", fmt.Errorf("указан неверный формат repeat")
	}

	resDate := parseDate.AddDate(0, 0, number)

	for resDate.Before(now) {
		resDate = resDate.AddDate(0, 0, number)
	}

	return resDate.Format("20060102"), nil
}

func addWeek(now, parseDate time.Time, repeat string) (string, error) {

	repeatArr := strings.Fields(repeat)

	if len(repeatArr) < 2 {
		return "", fmt.Errorf("указан неверный формат repeat")
	}

	days := strings.Split(repeatArr[1], ",")
	daysWeek := make([]int, len(days))

	for _, day := range days {
		num, err := strconv.Atoi(day)

		if err != nil {
			return "", fmt.Errorf("указан неверный формат repeat")
		}

		if num < 1 || num > 7 {
			return "", fmt.Errorf("указан неверный формат repeat")
		}

		daysWeek = append(daysWeek, num)
	}

	resDate := parseDate

	if resDate.Before(now) {
		resDate = now.AddDate(0, 0, 1)
	}

	for {
		weekday := int(resDate.Weekday())

		resDate = resDate.AddDate(0, 0, 1)

		for _, day := range daysWeek {
			if weekday == day-1 {
				return resDate.Format("20060102"), nil
			}
		}
	}
}

func addMonth(now, parseDate time.Time, repeat string) (string, error) {

	repeatArr := strings.Fields(repeat)

	if len(repeatArr) < 2 {
		return "", fmt.Errorf("указан неверный формат repeat")
	}

	resDate := parseDate

	if resDate.Before(now) {
		resDate = now
	}

	days, months, err := splitRepeatArr(repeatArr)

	if err != nil {
		return "", err
	}

	if len(repeatArr) < 3 {
		for {
			resDate = resDate.AddDate(0, 0, 1)

			for _, day := range days {
				if day == -1 && resDate.Day() == daysOfMonth(resDate).Day() {
					return resDate.Format("20060102"), nil
				} else if day == -2 && resDate.Day() == daysOfMonth(resDate).AddDate(0, 0, -1).Day() {
					return resDate.Format("20060102"), nil
				} else if resDate.Day() == day && day > 0 {
					return resDate.Format("20060102"), nil
				}
			}
		}
	} else {
		for {
			monthNum := int(resDate.Month())
			count := 0

			for _, month := range months {
				if monthNum == month {
					break
				}
				count++
			}
			if count == len(months) {
				resDate = resDate.AddDate(0, 1, 0)
				continue
			}
			if resDate.Equal(now) {
				resDate = resDate.AddDate(0, 1, 0)
				continue
			}

			for _, day := range days {
				if day == -1 {
					day = daysOfMonth(resDate).Day()
				} else if day == -2 {
					day = daysOfMonth(resDate).Day() - 1
				}

				result := time.Date(resDate.Year(), time.Month(monthNum), day, 0, 0, 0, 0, time.UTC)

				return result.Format("20060102"), nil
			}
		}
	}
}

func splitRepeatArr(repeatArr []string) ([]int, []int, error) {

	if len(repeatArr) > 3 {
		return nil, nil, fmt.Errorf("указан неверный формат repeat")
	}

	var days []int
	var months []int

	strDays := strings.Split(repeatArr[1], ",")

	for _, day := range strDays {
		num, err := strconv.Atoi(day)

		if err != nil || num < -2 || num > 31 {
			return nil, nil, fmt.Errorf("указан неверный формат repeat")
		}

		days = append(days, num)
	}

	if len(repeatArr) == 3 {
		strMonth := strings.Split(repeatArr[2], ",")

		for _, day := range strMonth {
			num, err := strconv.Atoi(day)

			if err != nil || num < 1 || num > 12 {
				return nil, nil, fmt.Errorf("указан неверный формат repeat")
			}

			months = append(months, num)
		}
	}

	return days, months, nil
}

func daysOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)
}

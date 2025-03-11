package repeattask

import "time"

func addYear(now, parseDate time.Time) string {
	resDate := parseDate

	if now.Before(resDate) {
		resDate = resDate.AddDate(1, 0, 0)
	} else {
		for !now.Before(resDate) {
			resDate = resDate.AddDate(1, 0, 0)
		}
	}

	return resDate.Format("20060102")
}

func addDay(now, parseDate time.Time, number int) string {
	resDate := parseDate

	if now.Before(resDate) {
		resDate = resDate.AddDate(0, 0, number)
	} else {
		for !now.Before(resDate) {
			resDate = resDate.AddDate(0, 0, number)
		}
	}

	return resDate.Format("20060102")
}

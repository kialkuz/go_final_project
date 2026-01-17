package nextDate

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"

	datePkg "github.com/Yandex-Practicum/final/pkg/date"
)

const (
	intervalTypeDays      = "d"
	intervalTypeWeekDays  = "w"
	intervalTypeMonthDays = "m"
	intervalTypeYear      = "y"
)

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	formatParts := strings.Split(repeat, " ")

	err := checkDateRepeat(now, formatParts)
	if err != nil {
		return "", err
	}

	nextDate, err := getNextDate(now, dstart, formatParts)
	if err != nil {
		return "", err
	}

	return nextDate.Format(datePkg.DateFormat), nil
}

func checkDateRepeat(now time.Time, formatParts []string) error {
	allowedIntervalTypes := []string{intervalTypeDays, intervalTypeWeekDays, intervalTypeMonthDays, intervalTypeYear}
	if !slices.Contains(allowedIntervalTypes, formatParts[0]) {
		return errors.New("недопустимый символ")
	}

	if len(formatParts) == 0 {
		return errors.New("отсутствует интервал")
	}

	var err error

	switch formatParts[0] {
	case intervalTypeDays:
		err = checkIntervalTypeDays(formatParts)
	case intervalTypeWeekDays:
		err = CheckIntervalTypeWeekDays(formatParts)
	case intervalTypeMonthDays:
		err = checkIntervalTypeMonthDays(now, formatParts)
	}

	return err
}

func getNextDate(now time.Time, dstart string, formatParts []string) (time.Time, error) {
	date, err := time.Parse(datePkg.DateFormat, dstart)

	var taskNextDate time.Time

	if err != nil {
		return taskNextDate, err
	}

	switch formatParts[0] {
	case intervalTypeDays:
		days, _ := strconv.Atoi(formatParts[1])

		taskNextDate = getNextDateByInterval(now, date, 0, days)
	case intervalTypeYear:
		taskNextDate = getNextDateByInterval(now, date, 1, 0)
	case intervalTypeWeekDays:
		taskNextDate = GetNextDateByWeekDays(now, date, formatParts[1])
	case intervalTypeMonthDays:
		var monthes string
		if len(formatParts) == 3 {
			monthes = formatParts[2]
		}

		taskNextDate = GetNextDateByMonthDays(now, date, formatParts[1], monthes)
	}

	return taskNextDate, nil
}

func getMonthLastDay(year, month int) int {
	return getMonthDay(year, month+1, 0).Day()
}

func getMonthDay(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func getNextDateByInterval(now, date time.Time, year, days int) time.Time {
	for {
		date = date.AddDate(year, 0, days)
		if datePkg.IsDateAfter(date, now) {
			return date
		}
	}
}

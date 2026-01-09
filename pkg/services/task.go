package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/final/pkg/db"
	"github.com/Yandex-Practicum/final/settings"
)

const (
	intervalTypeDays      = "d"
	intervalTypeWeekDays  = "w"
	intervalTypeMonthDays = "m"
	intervalTypeYear      = "y"

	countWeekDays = 7
	dateFormat    = "20060102"
)

var task db.Task

func GetTaskBody(r *http.Request) (*db.Task, error) {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return nil, errors.New("ошибка чтения тела запроса")
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		return nil, errors.New("ошибка десериализации JSON")
	}

	return &task, nil
}

func CheckTask(task *db.Task) (*db.Task, error) {
	if task.Title == "" {
		return nil, errors.New("отсутствует заголовок задачи")
	}

	now := time.Now()

	if task.Date != "" {
		t, err := time.Parse("20060102", task.Date)
		if err != nil {
			return nil, errors.New("дата представлена в формате, отличном от 20060102")
		}

		if isDateAfter(now, t) {
			if len(task.Repeat) == 0 {
				// если правила повторения нет, то берём сегодняшнее число
				task.Date = now.Format("20060102")
			} else {
				// в противном случае, берём вычисленную ранее следующую дату
				next, err := NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return nil, errors.New("ошибка получения даты задачи")
				}
				task.Date = next
			}
		}
	} else {
		task.Date = now.Format("20060102")
	}

	return task, nil
}

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

	return nextDate.Format(dateFormat), nil
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
		err = checkIntervalTypeWeekDays(formatParts)
	case intervalTypeMonthDays:
		err = checkIntervalTypeMonthDays(now, formatParts)
	}

	return err
}

func checkIntervalTypeDays(formatParts []string) error {
	if len(formatParts) < 2 || formatParts[1] == "" {
		return errors.New("неверный формат интервала")
	}

	count, err := strconv.Atoi(formatParts[1])
	if err != nil {
		return errors.New("количество должно быть цифрой")
	}

	if count <= 0 {
		return errors.New("передан нулевой или отрицательный интервал")
	}

	if count > settings.ServerSettings.MaxIntervalDays {
		return errors.New("превышен максимально допустимый интервал")
	}

	return nil
}

func checkIntervalTypeWeekDays(formatParts []string) error {
	if len(formatParts) < 2 || formatParts[1] == "" {
		return errors.New("неверный формат интервала")
	}

	weekDays := strings.Split(formatParts[1], ",")
	for _, weekDay := range weekDays {
		weekDay, err := strconv.Atoi(weekDay)
		if err != nil {
			return errors.New("неверный формат интервала")
		}
		if weekDay <= 0 {
			return errors.New("передан нулевой или отрицательный интервал")
		}

		if weekDay > countWeekDays {
			return errors.New("недопустимый день недели")
		}
	}

	return nil
}

func checkIntervalTypeMonthDays(now time.Time, formatParts []string) error {
	monthDays := strings.Split(formatParts[1], ",")
	lastMonthDay := getMonthLastDay(now.Year(), int(now.Month()))
	hasMonthesList := len(formatParts) == 3

	for _, monthDay := range monthDays {
		monthDayNumber, err := strconv.Atoi(monthDay)
		if err != nil {
			return errors.New("неверный формат интервала")
		}
		if monthDayNumber < -2 || (!hasMonthesList && monthDayNumber > lastMonthDay) {
			return errors.New("недопустимый день месяца")
		}
	}

	if hasMonthesList {
		monthes := strings.Split(formatParts[2], ",")
		decemberNumber := int(time.December)

		for _, month := range monthes {
			monthNumber, err := strconv.Atoi(month)
			if err != nil {
				return errors.New("неверный формат списка месяцев")
			}

			if monthNumber < 1 || monthNumber > decemberNumber {
				return errors.New("недопустимый месяц")
			}

			lastMonthDay := getMonthLastDay(now.Year(), monthNumber)

			for _, monthDay := range monthDays {
				monthDayNumber, _ := strconv.Atoi(monthDay)
				if monthDayNumber > lastMonthDay {
					return errors.New("недопустимый день месяца")
				}
			}
		}
	}

	return nil
}

func getNextDate(now time.Time, dstart string, formatParts []string) (time.Time, error) {
	date, err := time.Parse("20060102", dstart)

	var nextDate time.Time

	if err != nil {
		return nextDate, err
	}

	switch formatParts[0] {
	case intervalTypeDays:
		days, _ := strconv.Atoi(formatParts[1])

		nextDate = getNextDateByInterval(now, date, 0, days)
	case intervalTypeYear:
		nextDate = getNextDateByInterval(now, date, 1, 0)
	case intervalTypeWeekDays:
		nextDate = getNextDateByWeekDays(now, date, formatParts[1])
	case intervalTypeMonthDays:
		var monthes string
		if len(formatParts) == 3 {
			monthes = formatParts[2]
		}

		nextDate = getNextDateByMonthDays(now, date, formatParts[1], monthes)
	}

	return nextDate, nil
}

func getNextDateByWeekDays(now, date time.Time, days string) time.Time {
	weekDays := strings.Split(days, ",")

	var weekDaysNumbers []int

	for _, weekDay := range weekDays {
		monthDayNumber, _ := strconv.Atoi(weekDay)
		weekDaysNumbers = append(weekDaysNumbers, monthDayNumber)
	}

	currentWeekdayNumber := int(now.Weekday())

	var nextDate time.Time
	nextDateFound := false

	if currentWeekdayNumber == 0 {
		nextDate = getNextDateByInterval(now, date, 0, weekDaysNumbers[0])
		nextDateFound = true
	}

	if !nextDateFound {
		nextDate, nextDateFound = getNextDateByWeekDaysList(now, date, currentWeekdayNumber, weekDaysNumbers)
	}

	if !nextDateFound {
		nextDate = getNextDateByFirstListDay(weekDaysNumbers[0], now, date)
	}

	return nextDate
}

func getNextDateByWeekDaysList(now, date time.Time, currentWeekdayNumber int, weekDays []int) (time.Time, bool) {
	var nextDate time.Time
	nextDateFound := false

	for _, day := range weekDays {
		if day > currentWeekdayNumber {
			nextDate = getNextDateByInterval(now, date, 0, day-currentWeekdayNumber)
			nextDateFound = true
		}
	}

	return nextDate, nextDateFound
}

func getNextDateByFirstListDay(firstListDay int, now, date time.Time) time.Time {
	for {
		date = date.AddDate(0, 0, 1)
		if int(date.Weekday()) == firstListDay && isDateAfter(date, now) {
			break
		}
	}

	return date
}

func getNextDateByMonthDays(now, date time.Time, days, monthes string) time.Time {
	monthDays := strings.Split(days, ",")

	var monthDaysNumbers []int

	for _, monthDay := range monthDays {
		monthDayNumber, _ := strconv.Atoi(monthDay)
		monthDaysNumbers = append(monthDaysNumbers, monthDayNumber)
	}

	monthDaysNumbers = sortMonthDaysNumbers(monthDaysNumbers)

	if monthes == "" {
		dateStart := getDateStartForEveryMonth(now, date)
		nextDate, nextDateFound := getNextDateSelectedMonth(now, dateStart, monthDaysNumbers)
		if !nextDateFound {
			for {
				nextMonthDate := dateStart.AddDate(0, 1, 0)
				lastMonthDay := getMonthLastDay(nextMonthDate.Year(), int(nextMonthDate.Month()))

				if lastMonthDay >= monthDaysNumbers[0] {
					nextDate = getMonthDay(nextMonthDate.Year(), int(nextMonthDate.Month()), monthDaysNumbers[0])
					break
				}
			}
		}

		return nextDate
	} else {
		var monthNumbers []int

		monthList := strings.Split(monthes, ",")
		for _, month := range monthList {
			monthNumber, _ := strconv.Atoi(month)
			monthNumbers = append(monthNumbers, monthNumber)
		}

		slices.Sort(monthNumbers)

		dateStart := getDateStartForMonthesList(now, date, monthNumbers)
		nextDate, _ := getNextDateSelectedMonth(now, dateStart, monthDaysNumbers)

		return nextDate
	}
}

func sortMonthDaysNumbers(monthDaysNumbers []int) []int {
	var positiveNumbers []int
	var negativeNumbers []int

	for _, monthDaysNumber := range monthDaysNumbers {
		if monthDaysNumber > 0 {
			positiveNumbers = append(positiveNumbers, monthDaysNumber)
		} else {
			negativeNumbers = append(negativeNumbers, monthDaysNumber)
		}
	}

	slices.Sort(positiveNumbers)
	slices.Sort(negativeNumbers)

	return append(positiveNumbers, negativeNumbers...)
}

func getDateStartForEveryMonth(now, dateStart time.Time) time.Time {
	if isDateAfter(dateStart, now) {
		return dateStart
	}

	return now
}

func getDateStartForMonthesList(now, dateStart time.Time, monthes []int) time.Time {
	currentMonthNumber := int(now.Month())
	dateStartMonthNumber := int(dateStart.Month())
	isDateAfter := isDateAfter(dateStart, now)

	for _, monthNumber := range monthes {
		if currentMonthNumber == monthNumber && dateStartMonthNumber == monthNumber && isDateAfter {
			return dateStart
		} else if currentMonthNumber < monthNumber {
			return getMonthDay(now.Year(), monthNumber, 1)
		}
	}

	return getMonthDay(now.Year()+1, monthes[0], 1)
}

func getNextDateSelectedMonth(now, date time.Time, monthDaysNumbers []int) (time.Time, bool) {
	nextDate := date

	dateStartMonthDay := date.Day()
	nextDateFound := false

	lastMonthDay := getMonthLastDay(date.Year(), int(date.Month()))

	for _, dayNumber := range monthDaysNumbers {
		if lastMonthDay < dayNumber {
			continue
		}

		if dateStartMonthDay < dayNumber {
			nextDate = getNextDateByInterval(now, date, 0, dayNumber-dateStartMonthDay)
			nextDateFound = true
			break
		} else if dayNumber == -2 || dayNumber == -1 {
			prevMonthDay := lastMonthDay - 1

			if dayNumber == -2 && dateStartMonthDay < prevMonthDay {
				nextDate = getMonthDay(date.Year(), int(date.Month()), prevMonthDay)
				nextDateFound = true
				break
			} else if dayNumber == -1 && dateStartMonthDay < lastMonthDay {
				nextDate = getMonthDay(date.Year(), int(date.Month()), lastMonthDay)
				nextDateFound = true
				break
			}
		}
	}

	return nextDate, nextDateFound
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
		if isDateAfter(date, now) {
			return date
		}
	}
}

func isDateAfter(firstDate, secondDate time.Time) bool {
	return firstDate.Format(dateFormat) > secondDate.Format(dateFormat)
}

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

		if isAfterNow(now, t) {
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

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	formatParts := strings.Split(repeat, " ")

	err := checkDateRepeat(formatParts)
	if err != nil {
		return "", err
	}

	nextDate, err := getNextDate(now, dstart, formatParts)
	if err != nil {
		return "", err
	}

	return nextDate.Format(dateFormat), nil
}

func checkDateRepeat(formatParts []string) error {
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
		err = checkIntervalTypMonthDays(formatParts)
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

func checkIntervalTypMonthDays(formatParts []string) error {
	now := time.Now()

	monthDays := strings.Split(formatParts[1], ",")
	lastMonthDay := getLastMonthDay(now.Year(), int(now.Month())+1)

	for monthDay := range monthDays {
		if monthDay <= 0 || monthDay > lastMonthDay || monthDay != -1 || monthDay != -2 {
			return errors.New("недопустимый день месяца")
		}
	}

	if len(formatParts) == 3 {
		monthes := strings.Split(formatParts[2], ",")
		decemberNumber := int(time.December)

		for month := range monthes {
			if month <= 0 || month > decemberNumber {
				return errors.New("недопустимый месяц")
			}
		}
	}

	return nil
}

func getLastMonthDay(year, month int) int {
	return time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
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
		// case intervalTypeWeekDays:
		// 	nextDate = checkIntervalTypeWeekDays(formatParts[1])
		/*case intervalTypeMonthDays:
		nextDate = checkIntervalTypMonthDays(formatParts)*/
	}

	return nextDate, nil
}

func getNextDateByInterval(now, date time.Time, year, days int) time.Time {
	for {
		date = date.AddDate(year, 0, days)
		if isAfterNow(date, now) {
			return date
		}
	}
}

func isAfterNow(date, now time.Time) bool {
	return date.Format(dateFormat) > now.Format(dateFormat)
}

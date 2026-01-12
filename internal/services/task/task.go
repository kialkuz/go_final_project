package task

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/final/internal/dto"
	"github.com/Yandex-Practicum/final/internal/services/task/nextDate"
	datePkg "github.com/Yandex-Practicum/final/pkg/date"
)

var task dto.Task

func GetTaskBody(r *http.Request) (*dto.Task, error) {
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

func CheckTask(task *dto.Task) error {
	if task.Title == "" {
		return errors.New("отсутствует заголовок задачи")
	}

	if task.Date != "" {
		_, err := time.Parse(datePkg.DateFormat, task.Date)
		if err != nil {
			return errors.New("дата представлена в формате, отличном от 20060102")
		}
	}

	return nil
}

func GetDateByRules(task *dto.Task) (string, error) {
	now := time.Now()

	var dateByRules string
	if task.Date != "" {
		t, _ := time.Parse(datePkg.DateFormat, task.Date)
		dateByRules = t.Format(datePkg.DateFormat)

		if datePkg.IsDateAfter(now, t) {
			if len(task.Repeat) == 0 {
				// если правила повторения нет, то берём сегодняшнее число
				dateByRules = now.Format(datePkg.DateFormat)
			} else {
				// в противном случае, берём вычисленную ранее следующую дату
				next, err := nextDate.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return "", errors.New("ошибка получения даты задачи")
				}
				dateByRules = next
			}
		}
	} else {
		dateByRules = now.Format(datePkg.DateFormat)
	}

	return dateByRules, nil
}

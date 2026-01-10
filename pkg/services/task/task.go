package task

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/final/pkg/dto"
	dateService "github.com/Yandex-Practicum/final/pkg/services/date"
	"github.com/Yandex-Practicum/final/pkg/services/task/nextdate"
)

const (
	maxIntervalDaysEnv     = "MAX_INTERNAL_DAYS"
	defaultMaxIntervalDays = "400"
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
		_, err := time.Parse("20060102", task.Date)
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
		t, _ := time.Parse("20060102", task.Date)
		dateByRules = t.Format(dateService.DateFormat)

		if dateService.IsDateAfter(now, t) {
			if len(task.Repeat) == 0 {
				// если правила повторения нет, то берём сегодняшнее число
				dateByRules = now.Format("20060102")
			} else {
				// в противном случае, берём вычисленную ранее следующую дату
				next, err := nextdate.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					return "", errors.New("ошибка получения даты задачи")
				}
				dateByRules = next
			}
		}
	} else {
		dateByRules = now.Format("20060102")
	}

	return dateByRules, nil
}

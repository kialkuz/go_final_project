package task

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Yandex-Practicum/final/internal/dto"
	"github.com/Yandex-Practicum/final/internal/infrastructure/repository"
	httpService "github.com/Yandex-Practicum/final/internal/services/http"
	"github.com/Yandex-Practicum/final/internal/services/task/nextDate"
)

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Println("Не передан ID " + id)
		httpService.WriteJsonBadRequest(w, dto.ErrorResponse{ErrorText: "Не передан ID"})
		return
	}

	numericId, err := strconv.Atoi(id)
	if err != nil {
		log.Println("Передан некорректный ID " + id)
		httpService.WriteJsonBadRequest(w, dto.ErrorResponse{ErrorText: "Передан некорректный ID"})
		return
	}

	task, err := repository.GetTask(numericId)
	if err != nil {
		log.Println(err)
		httpService.WriteJson(w, dto.ErrorResponse{ErrorText: "Задача не найдена"}, http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		err = repository.DeleteTask(numericId)
		if err != nil {
			log.Println(err)
			httpService.WriteJsonInternalServerError(w, dto.ErrorResponse{ErrorText: "Ошибка удаления задачи"})
			return
		}
	} else {
		nextDate, err := nextDate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			log.Println(err)
			httpService.WriteJsonInternalServerError(
				w,
				dto.ErrorResponse{ErrorText: "Ошибка отметки задачи выполненной"},
			)
			return
		}

		task.Date = nextDate
		err = repository.UpdateTask(task)
		if err != nil {
			log.Println(err)
			httpService.WriteJsonInternalServerError(
				w,
				dto.ErrorResponse{ErrorText: "Ошибка отметки задачи выполненной"},
			)
			return
		}
	}

	httpService.WriteJsonOKResponse(w, dto.EmptyResponse{})
}

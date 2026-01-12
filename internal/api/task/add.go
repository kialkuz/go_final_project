package task

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Yandex-Practicum/final/internal/dto"
	"github.com/Yandex-Practicum/final/internal/infrastructure/repository"
	httpService "github.com/Yandex-Practicum/final/internal/services/http"
	taskService "github.com/Yandex-Practicum/final/internal/services/task"
)

type AddedTask struct {
	ID string `json:"id"`
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := taskService.GetTaskBody(r)

	if err != nil {
		log.Println(err)
		httpService.WriteJsonBadRequest(w, dto.ErrorResponse{ErrorText: err.Error()})
	}

	err = taskService.CheckTask(task)
	if err != nil {
		log.Println(err)
		httpService.WriteJsonBadRequest(w, dto.ErrorResponse{ErrorText: err.Error()})
		return
	}

	task.Date, err = taskService.GetDateByRules(task)
	if err != nil {
		log.Println(err)
		httpService.WriteJsonInternalServerError(w, dto.ErrorResponse{ErrorText: err.Error()})
		return
	}

	id, err := repository.AddTask(task)

	if err != nil {
		log.Println(err)
		httpService.WriteJsonInternalServerError(w, dto.ErrorResponse{ErrorText: "ошибка добавления задачи"})
		return
	}

	httpService.WriteJson(w, AddedTask{
		ID: strconv.FormatInt(id, 10),
	}, http.StatusCreated)
}

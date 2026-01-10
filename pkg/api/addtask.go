package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Yandex-Practicum/final/pkg/dto"
	"github.com/Yandex-Practicum/final/pkg/infrastructure/repository"
	taskService "github.com/Yandex-Practicum/final/pkg/services/task"
)

type AddedTask struct {
	ID string `json:"id"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := taskService.GetTaskBody(r)
	if err != nil {
		log.Println(err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: err.Error()})
	}

	err = taskService.CheckTask(task)
	if err != nil {
		log.Println(err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: err.Error()})
		return
	}

	task.Date, err = taskService.GetDateByRules(task)
	if err != nil {
		log.Println(err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: err.Error()})
		return
	}

	id, err := repository.AddTask(task)

	if err != nil {
		log.Println(err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: "ошибка добавления задачи"})
		return
	}

	writeJson(w, AddedTask{
		ID: strconv.FormatInt(id, 10),
	})
}

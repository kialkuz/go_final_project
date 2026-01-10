package api

import (
	"log"
	"net/http"

	"github.com/Yandex-Practicum/final/pkg/dto"
	"github.com/Yandex-Practicum/final/pkg/infrastructure/repository"
	taskService "github.com/Yandex-Practicum/final/pkg/services/task"
)

type EditTask struct {
	ID string `json:"id"`
}

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := taskService.GetTaskBody(r)
	if err != nil {
		log.Println(err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: err.Error()})
		return
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

	err = repository.UpdateTask(task)
	if err != nil {
		log.Println(err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: "ошибка редактирования задачи"})
		return
	}

	writeJson(w, AddedTask{
		ID: task.ID,
	})
}

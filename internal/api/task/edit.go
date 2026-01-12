package task

import (
	"log"
	"net/http"

	"github.com/Yandex-Practicum/final/internal/dto"
	"github.com/Yandex-Practicum/final/internal/infrastructure/repository"
	httpService "github.com/Yandex-Practicum/final/internal/services/http"
	taskService "github.com/Yandex-Practicum/final/internal/services/task"
)

type EditTask struct {
	ID string `json:"id"`
}

func EditTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := taskService.GetTaskBody(r)
	if err != nil {
		log.Println(err)
		httpService.WriteJsonBadRequest(w, dto.ErrorResponse{ErrorText: err.Error()})
		return
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

	err = repository.UpdateTask(task)
	if err != nil {
		log.Println(err)
		httpService.WriteJsonInternalServerError(w, dto.ErrorResponse{ErrorText: "ошибка редактирования задачи"})
		return
	}

	httpService.WriteJsonOKResponse(w, AddedTask{
		ID: task.ID,
	})
}

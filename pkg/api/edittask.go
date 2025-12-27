package api

import (
	"log"
	"net/http"

	"github.com/Yandex-Practicum/final/pkg/db"
	"github.com/Yandex-Practicum/final/pkg/dto"
	"github.com/Yandex-Practicum/final/pkg/services"
)

type EditTask struct {
	ID string `json:"id"`
}

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := services.GetTaskBody(r)
	if err != nil {
		log.Println(err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: err.Error()})
		return
	}

	task, err = services.CheckTask(task)
	if err != nil {
		log.Println(err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: err.Error()})
		return
	}

	err = db.UpdateTask(task)
	if err != nil {
		log.Println(err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: "ошибка редактирования задачи"})
		return
	}

	writeJson(w, AddedTask{
		ID: task.ID,
	})
}

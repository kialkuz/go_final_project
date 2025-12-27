package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Yandex-Practicum/final/pkg/db"
	"github.com/Yandex-Practicum/final/pkg/dto"
	"github.com/Yandex-Practicum/final/pkg/services"
)

type AddedTask struct {
	ID string `json:"id"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := services.GetTaskBody(r)
	if err != nil {
		log.Println(err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: err.Error()})
	}

	task, err = services.CheckTask(task)
	if err != nil {
		log.Println(err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: err.Error()})
		return
	}

	id, err := db.AddTask(task)

	if err != nil {
		log.Println(err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: "ошибка добавления задачи"})
		return
	}

	writeJson(w, AddedTask{
		ID: strconv.FormatInt(id, 10),
	})
}

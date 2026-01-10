package api

import (
	"log"
	"net/http"

	"github.com/Yandex-Practicum/final/pkg/dto"
	"github.com/Yandex-Practicum/final/pkg/infrastructure/repository"
)

type TasksResp struct {
	Tasks []*dto.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := repository.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		log.Println(err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: "ошибка получения"})
		return
	}

	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

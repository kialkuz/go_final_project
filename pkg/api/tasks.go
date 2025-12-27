package api

import (
	"log"
	"net/http"

	"github.com/Yandex-Practicum/final/pkg/db"
	"github.com/Yandex-Practicum/final/pkg/dto"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		log.Println(err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: "ошибка получения"})
		return
	}

	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

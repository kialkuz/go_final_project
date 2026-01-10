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
	var tasks []*dto.Task
	var err error

	q := r.URL.Query()

	search := q.Get("search")
	if search == "" {
		tasks, err = repository.Tasks(50) // в параметре максимальное количество записей
	} else {
		tasks, err = repository.SearchTasks(search, 50)
	}

	if err != nil {
		log.Println(err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: "ошибка получения"})
		return
	}

	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

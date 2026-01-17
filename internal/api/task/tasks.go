package task

import (
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/final/internal/dto"
	"github.com/Yandex-Practicum/final/internal/infrastructure/repository"
	httpService "github.com/Yandex-Practicum/final/internal/services/http"
)

type TasksResp struct {
	Tasks []*dto.Task `json:"tasks"`
}

const tasksLimit = 50

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []*dto.Task
	var err error

	q := r.URL.Query()

	search := q.Get("search")
	if search == "" {
		tasks, err = repository.Tasks(50) // в параметре максимальное количество записей
	} else {
		var searchDate time.Time
		searchDate, err = time.Parse("02.01.2006", search)

		if err != nil {
			tasks, err = repository.SearchTasksByText(search, tasksLimit)
		} else {
			tasks, err = repository.SearchTasksByDate(searchDate, tasksLimit)
		}
	}

	if err != nil {
		log.Println(err)
		httpService.WriteJsonInternalServerError(w, dto.ErrorResponse{ErrorText: "ошибка получения"})
		return
	}

	httpService.WriteJsonOKResponse(w, TasksResp{
		Tasks: tasks,
	})
}

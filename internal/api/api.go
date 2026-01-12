package api

import (
	"net/http"

	"github.com/Yandex-Practicum/final/internal/api/task"
	"github.com/Yandex-Practicum/final/internal/infrastructure/env"
	"github.com/Yandex-Practicum/final/internal/middleware"
)

func Init() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /", http.FileServer(http.Dir(env.EnvList.WebDirPath)))
	mux.HandleFunc("GET /api/nextdate", task.NextDayHandler)
	mux.HandleFunc("GET /api/task", middleware.Auth(task.TaskHandler))
	mux.HandleFunc("POST /api/task", middleware.Auth(task.AddTaskHandler))
	mux.HandleFunc("PUT /api/task", middleware.Auth(task.EditTaskHandler))
	mux.HandleFunc("GET /api/tasks", middleware.Auth(task.TasksHandler))
	mux.HandleFunc("DELETE /api/task", task.DeleteTaskHandler)
	mux.HandleFunc("POST /api/task/done", middleware.Auth(task.DoneTaskHandler))
	mux.HandleFunc("POST /api/signin", signinHandler)

	return mux
}

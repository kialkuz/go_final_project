package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Yandex-Practicum/final/pkg/dto"
	"github.com/Yandex-Practicum/final/pkg/infrastructure/env"
	"github.com/Yandex-Practicum/final/pkg/middleware"
)

const (
	webDirPathEnv     = "WEB_DIR_PATH"
	defaultWebDirPath = "7540"
)

func Init() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServer(http.Dir(env.Lookup(webDirPathEnv, defaultWebDirPath))))
	mux.HandleFunc("GET /api/nextdate", nextDayHandler)
	mux.HandleFunc("GET /api/task", middleware.Auth(taskHandler))
	mux.HandleFunc("POST /api/task", middleware.Auth(addTaskHandler))
	mux.HandleFunc("PUT /api/task", middleware.Auth(editTaskHandler))
	mux.HandleFunc("GET /api/tasks", middleware.Auth(tasksHandler))
	mux.HandleFunc("DELETE /api/task", deleteTaskHandler)
	mux.HandleFunc("POST /api/task/done", middleware.Auth(doneTaskHandler))
	mux.HandleFunc("POST /api/signin", signinHandler)

	return mux
}

func writeJson(w http.ResponseWriter, data any) {
	resp, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, ok := data.(dto.ErrorResponse)
	if ok {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Write(resp)
}

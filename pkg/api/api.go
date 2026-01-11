package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Yandex-Practicum/final/pkg/dto"
	"github.com/Yandex-Practicum/final/pkg/infrastructure/env"
)

const (
	webDirPathEnv     = "WEB_DIR_PATH"
	defaultWebDirPath = "7540"
)

func Init() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServer(http.Dir(env.Lookup(webDirPathEnv, defaultWebDirPath))))
	mux.HandleFunc("GET /api/nextdate", nextDayHandler)
	mux.HandleFunc("POST /api/task", addTaskHandler)
	mux.HandleFunc("GET /api/tasks", tasksHandler)
	mux.HandleFunc("GET /api/task", taskHandler)
	mux.HandleFunc("PUT /api/task", editTaskHandler)
	mux.HandleFunc("DELETE /api/task", deleteTaskHandler)
	mux.HandleFunc("POST /api/task/done", doneTaskHandler)

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

package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Yandex-Practicum/final/pkg/dto"
	"github.com/Yandex-Practicum/final/pkg/infrastructure/repository"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Println("Не передан ID " + id)
		writeJson(w, dto.ErrorResponse{ErrorText: "Не передан ID"})
		return
	}

	numericId, err := strconv.Atoi(id)
	if err != nil {
		log.Println("Передан некорректный ID " + id)
		writeJson(w, dto.ErrorResponse{ErrorText: "Передан некорректный ID"})
		return
	}

	task, err := repository.GetTask(numericId)
	if err != nil {
		log.Println(err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: "Задача не найдена"})
		return
	}

	writeJson(w, task)
}

package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Yandex-Practicum/final/pkg/dto"
	"github.com/Yandex-Practicum/final/pkg/infrastructure/repository"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	_, err = repository.GetTask(numericId)
	if err != nil {
		log.Println(err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: "Задача не найдена"})
		return
	}

	err = repository.DeleteTask(numericId)
	if err != nil {
		log.Println(err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: "Ошибка удаления задачи"})
		return
	}

	responce, _ := json.Marshal(dto.EmptyResponse{})
	w.Write([]byte(responce))
}

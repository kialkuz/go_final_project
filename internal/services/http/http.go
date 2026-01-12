package services

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Yandex-Practicum/final/internal/dto"
)

func WriteJsonOKResponse(w http.ResponseWriter, data any) {
	resp := serializeJson(w, data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func WriteJsonBadRequest(w http.ResponseWriter, data any) {
	resp := serializeJson(w, data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(resp)
}

func WriteJsonInternalServerError(w http.ResponseWriter, data any) {
	resp := serializeJson(w, data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(resp)
}

func WriteJson(w http.ResponseWriter, data any, statusCode int) {
	resp := serializeJson(w, data)

	w.Header().Set("Content-Type", "application/json")

	_, ok := data.(dto.ErrorResponse)
	if ok {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Write(resp)
}

func serializeJson(w http.ResponseWriter, data any) []byte {
	resp, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return resp
}

func WriteJsonWithoutSerialize(w http.ResponseWriter, resp []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resp)
}

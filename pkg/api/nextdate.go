package api

import (
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/final/pkg/dto"
	"github.com/Yandex-Practicum/final/pkg/services"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	now, err := time.Parse("20060102", q.Get("now"))

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "error now date", http.StatusBadRequest)
		return
	}

	nextDate, err := services.NextDate(now, q.Get("date"), q.Get("repeat"))
	if err != nil {
		log.Println(err.Error())
		writeJson(w, dto.ErrorResponse{ErrorText: err.Error()})
	}

	w.Write([]byte(nextDate))
}

package api

import (
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/final/settings"
)

func Init() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServer(http.Dir(settings.ServerSettings.WebDir)))
	mux.HandleFunc("GET /api/nextdate", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		now, err := time.Parse("20060102", q.Get("now"))

		if err != nil {
			log.Println(err.Error())
			http.Error(w, "error now date", http.StatusBadRequest)
			return
		}

		nextDate, err := NextDate(now, q.Get("date"), q.Get("repeat"))
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "error get next date", http.StatusBadRequest)
			return
		}

		w.Write([]byte(nextDate))
	})

	return mux
}

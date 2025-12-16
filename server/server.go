package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Yandex-Practicum/final/settings"
)

func Handle(logger *log.Logger) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(settings.ServerSettings.WebDir)))

	Port := os.Getenv("TODO_PORT")
	if Port == "" {
		Port = settings.ServerSettings.DefaultPort
	}

	return &http.Server{
		Addr:         ":" + Port,
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  time.Duration(5 * time.Second),
		WriteTimeout: time.Duration(10 * time.Second),
		IdleTimeout:  time.Duration(15 * time.Second),
	}
}

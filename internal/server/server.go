package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/final/internal/api"
	"github.com/Yandex-Practicum/final/internal/infrastructure/env"
)

func Handle(logger *log.Logger) *http.Server {
	mux := api.Init()

	return &http.Server{
		Addr:         ":" + env.EnvList.Port,
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  time.Duration(5 * time.Second),
		WriteTimeout: time.Duration(10 * time.Second),
		IdleTimeout:  time.Duration(15 * time.Second),
	}
}

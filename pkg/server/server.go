package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/final/pkg/api"
	"github.com/Yandex-Practicum/final/pkg/infrastructure/env"
)

const (
	portEnv     = "PORT"
	defaultPort = "7540"
)

func Handle(logger *log.Logger) *http.Server {
	mux := api.Init()

	Port := env.Lookup(portEnv, defaultPort)

	return &http.Server{
		Addr:         ":" + Port,
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  time.Duration(5 * time.Second),
		WriteTimeout: time.Duration(10 * time.Second),
		IdleTimeout:  time.Duration(15 * time.Second),
	}
}

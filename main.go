package main

import (
	"log"

	"github.com/Yandex-Practicum/final/internal/bootstrap"
	"github.com/Yandex-Practicum/final/internal/server"
)

func main() {
	bootstrap.Init()
	defer bootstrap.Db.Close()

	server := server.Handle(log.Default())

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err.Error())
	}
}

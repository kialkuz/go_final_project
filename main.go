package main

import (
	"fmt"
	"log"

	"github.com/Yandex-Practicum/final/internal/bootstrap"
	"github.com/Yandex-Practicum/final/internal/infrastructure/env"
	"github.com/Yandex-Practicum/final/internal/server"
)

func main() {
	bootstrap.Init()
	defer bootstrap.Db.Close()

	server := server.Handle(log.Default())
	fmt.Println("Port for start: " + env.EnvList.Port)

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err.Error())
	}
}

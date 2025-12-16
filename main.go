package main

import (
	"log"
	"os"

	"github.com/Yandex-Practicum/final/pkg/db"
	"github.com/Yandex-Practicum/final/server"
	"github.com/Yandex-Practicum/final/settings"
)

func main() {
	settings.NewConfiguration()

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = settings.ServerSettings.DbFile
	}
	err := db.Init(dbFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	server := server.Handle(log.Default())

	err = server.ListenAndServe()

	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

}

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Yandex-Practicum/final/pkg/db"
	"github.com/Yandex-Practicum/final/pkg/server"
	"github.com/Yandex-Practicum/final/settings"
)

var rootPath string

func main() {
	initRootPath()
	settings.NewConfiguration()

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = settings.ServerSettings.DbFile
	}

	database, err := db.Init(DBPath(dbFile))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer database.Close()

	server := server.Handle(log.Default())

	err = server.ListenAndServe()

	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}
}

func initRootPath() {
	path, err := os.Executable()
	if err != nil {
		log.Fatalf("os.Executable: %v", err)
	}

	rootPath = filepath.Dir(path)
}

func DBPath(dbFileName string) string {
	return filepath.Join(rootPath, dbFileName)
}

func WebPath(webDirName string) string {
	return filepath.Join(rootPath, webDirName)
}

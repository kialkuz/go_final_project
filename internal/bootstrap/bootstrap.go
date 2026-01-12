package bootstrap

import (
	"database/sql"
	"log"
	"os"

	"github.com/Yandex-Practicum/final/internal/infrastructure/env"
	"github.com/Yandex-Practicum/final/pkg/path"
	_ "modernc.org/sqlite"
)

const schema = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX tasks_date ON scheduler (date);`

var Db *sql.DB

func Init() {
	path.InitRootPath()
	env.Load()

	err := initDB(env.EnvList.DbFile)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func initDB(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	Db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err = Db.Exec(schema)

		if err != nil {
			return err
		}
	}

	return nil
}

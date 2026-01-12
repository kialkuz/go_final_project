package env

import (
	"encoding/json"
	"log"
	"reflect"

	"github.com/joho/godotenv"
)

type Env struct {
	WebDirPath      string `json:"TODO_WEB_DIR_PATH"`
	Port            string `json:"TODO_PORT"`
	DbFile          string `json:"TODO_DB_FILE"`
	MaxIntervalDays string `json:"TODO_MAX_INTERNAL_DAYS"`
	Password        string `json:"TODO_PASSWORD"`
}

var EnvList Env

func Load() {
	envMap, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	for _, v := range envMap {
		if v == "" {
			log.Fatal("fill all variables in .env file")
		}
	}

	jsonData, err := json.Marshal(envMap)
	if err != nil {
		log.Fatal("Error convert env to json")
	}

	err = json.Unmarshal(jsonData, &EnvList)

	if err != nil {
		log.Fatal("Error create env struct")
	}

	values := reflect.ValueOf(EnvList)

	for i := 0; i < values.NumField(); i++ {
		if values.Field(i).String() == "" {
			log.Fatal("fill all variables in .env file")
		}
	}
}

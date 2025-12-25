package settings

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	WebDir          string `yaml:"web_dir"`
	DefaultPort     string `yaml:"default_port"`
	DbFile          string `yaml:"db_file"`
	MaxIntervalDays int    `yaml:"max_interval_days"`
}

var ServerSettings Settings

func NewConfiguration() {
	yamlFile, err := os.ReadFile("settings.yaml")
	if err != nil {
		log.Fatal("ошибка при чтении файла настроек", err)
	}

	err = yaml.Unmarshal(yamlFile, &ServerSettings)
	if err != nil {
		log.Fatal("ошибка при десериализации:", err)
		return
	}
}

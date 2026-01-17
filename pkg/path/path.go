package path

import (
	"log"
	"os"
	"path/filepath"
)

var rootPath string

func InitRootPath() {
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

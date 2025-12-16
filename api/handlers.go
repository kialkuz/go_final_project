package api

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func ViewIndex(res http.ResponseWriter, req *http.Request) {
	path, err := filepath.Abs("web/index.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(res, "error view file", http.StatusBadRequest)
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Println(err.Error())
		http.Error(res, "error view file", http.StatusBadRequest)
		return
	}

	res.Write([]byte(data))
}

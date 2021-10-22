package config

import (
	"fmt"
	"net/http"

	"github.com/pagarme/marshals/labs/vicki/desafioGo/logger"
)

var (
	version = "development"
)

func VersionHandler() {
	http.HandleFunc("/version", getVersion)
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	logger.InfoLog.Println("Version: ", version)
	fmt.Fprintf(w, "Version: %s", version)
}

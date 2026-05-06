package main

import (
	"os"

	"tracklink/internal/app"
	applogger "tracklink/internal/platform/logger"
)

func main() {
	log := applogger.New("text")

	if err := app.Run(); err != nil {
		log.Error("app_run_failed", "error", err)
		os.Exit(1)
	}
}

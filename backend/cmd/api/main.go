package main

import (
	"log"
	"os"

	"tracklink/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Printf("app: %v", err)
		os.Exit(1)
	}
}

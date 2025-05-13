package main

import (
	"log"

	"test_go/internal/app"
)
func main() {

	app, err := app.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := app.Router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
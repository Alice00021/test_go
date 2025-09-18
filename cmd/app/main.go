package main

import (
	_ "test_go/docs" // Импорт сгенерированной документации Swagger
)

//// @title Author API
//// @version 1.0
//// @description Это пример API для управления авторами и книгами
//// @host localhost:8080
//// @BasePath /api
//func main() {
//	app, err := app.NewApp()
//	if err != nil {
//		log.Fatalf("Не удалось инициализировать приложение: %v", err)
//	}
//
//	app.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
//
//	if err := app.Router.Run(":8080"); err != nil {
//		log.Fatalf("Не удалось запустить сервер: %v", err)
//	}
//}

import (
	"test_go/config"
	"test_go/internal/app"

	"log"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}

package main

/* package main

import (
    "log"
    "test_go/config"
    "test_go/models"
    "test_go/routes"
    "test_go/migrations"

    "github.com/gin-gonic/gin"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Ошибка загрузки конфигурации: %v", err)
    }

    db, err := models.Init_DB(cfg)
    if err != nil {
        log.Fatalf("Ошибка инициализации базы данных: %v", err)
    } */

    /* // Автомиграция для создания таблиц
    if err := db.AutoMigrate(&models.Author{}); err != nil {
        log.Fatalf("Ошибка миграции базы данных: %v", err)
    }
 */
    /* if err := migrations.RunMigrations(db); err != nil {
		panic("migration failed: " + err.Error())
	}

    router := gin.Default()

    routes.SetupRoutes(router, db)
    
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Ошибка запуска сервера: %v", err)
    }
} */
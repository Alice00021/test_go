package app

import (
	
	"fmt"
	"test_go/config"
	"test_go/db"
	"test_go/internal/controller/http"
	"test_go/internal/repo/pg"
	"test_go/internal/service"
	/* "test_go/migrations" */
	"test_go/routes"

	"github.com/gin-gonic/gin"

)

type App struct{
	Router *gin.Engine
}

func NewApp() (*App, error) {

	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	dbConn, err := db.Init_DB(cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации базы данных: %w", err)
	}

	/* // Выполнение миграций
	if err := migrations.RunMigrations(dbConn); err != nil {
		return nil, fmt.Errorf("ошибка выполнения миграций: %w", err)
	}
 */
	// Инициализация репозиториев
	bookRepo := pg.NewBookRepo(dbConn)
	authorRepo := pg.NewAuthorRepo(dbConn)
	userRepo := pg.NewUserRepo(dbConn)
	// Инициализация сервисов
	bookSvc := service.NewBookService(bookRepo)
	authorSvc := service.NewAuthorService(authorRepo)
	userSvc := service.NewAuthService(userRepo)
	// Инициализация обработчиков
	bookHandler := http.NewBookHandler(bookSvc)
	authorHandler := http.NewAuthorHandler(authorSvc)
	authHandler := http.NewAuthHandler(userSvc)
	// Настройка маршрутов
	router := gin.Default()
	routes.SetUpRoutes(router, bookHandler, authorHandler, authHandler)

	return &App{Router: router}, nil
}
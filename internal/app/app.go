package app

import (
	"fmt"
	"test_go/config"
	"test_go/db"
	"test_go/internal/controller"
	"test_go/internal/controller/http"
	"test_go/internal/repo/pg"
	"test_go/internal/service"
	"test_go/pkg/jwt"
	"test_go/routes"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router    *gin.Engine
	WsHub     *controller.Hub
	WsHandler *controller.WebSocketHandler
}

func NewApp() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	jwtManager := jwt.NewJWTManager(cfg.Jwt.SecretKey)
	dbConn, err := db.Init_DB(cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации базы данных: %w", err)
	}

	// Инициализация репозиториев
	bookRepo := pg.NewBookRepo(dbConn)
	authorRepo := pg.NewAuthorRepo(dbConn)
	userRepo := pg.NewUserRepo(dbConn)

	// Инициализация сервисов
	bookSvc := service.NewBookService(bookRepo)
	authorSvc := service.NewAuthorService(authorRepo)
	userSvc := service.NewAuthService(userRepo, jwtManager)

	// Инициализация WebSocket Hub
	wsHub := controller.NewHub()
	go wsHub.Run()

	// Инициализация обработчиков
	bookHandler := http.NewBookHandler(bookSvc)
	authorHandler := http.NewAuthorHandler(authorSvc, wsHub)
	authHandler := http.NewAuthHandler(userSvc)
	wsHandler := controller.NewWebSocketHandler(wsHub, userSvc)

	// Настройка маршрутов
	router := gin.Default()
	routes.SetUpRoutes(router, bookHandler, authorHandler, authHandler, jwtManager, wsHandler)

	return &App{
		Router:    router,
		WsHub:     wsHub,
		WsHandler: wsHandler,
	}, nil
}

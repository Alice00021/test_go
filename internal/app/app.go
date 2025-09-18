//package app
//
//import (
//	"fmt"
//	"sync"
//	"test_go/config"
//	"test_go/db"
//	"test_go/internal/controller/http/v1"
//	"test_go/internal/entity"
//	"test_go/internal/usecase/author"
//	"test_go/internal/usecase/book"
//	"test_go/internal/usecase/export"
//	"test_go/internal/usecase/user"
//	"test_go/pkg/jwt"
//	/* "test_go/migrations" */
//	"github.com/gin-gonic/gin"
//	"test_go/internal/di"
//	"test_go/routes"
//)
//
//type App struct {
//	Router *gin.Engine
//}
//
//func NewApp() (*App, error) {
//
//	cfg, err := config.Load()
//	if err != nil {
//		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
//	}
//	jwtManager := jwt.NewJWTManager(cfg.Jwt.SecretKey)
//
//	dbConn, err := db.Init_DB(cfg)
//	if err != nil {
//		return nil, fmt.Errorf("ошибка инициализации базы данных: %w", err)
//	}
//
//	// Repo
//	repo := di.NewRepo(pg)
//
//	// Use-Case
//	uc := di.NewUseCase(pgTx, repo, l, cfg)
//
//	txMtx := &sync.Mutex{}
//	// Инициализация репозиториев
//	bookRepo := pg.NewBookRepo(dbConn)
//	authorRepo := pg.NewAuthorRepo(dbConn)
//	userRepo := pg.NewUserRepo(dbConn)
//	// Инициализация сервисов
//	bookSvc := book.NewBookService(bookRepo)
//	authorSvc := author.NewAuthorService(authorRepo)
//	userSvc := user.NewAuthService(userRepo, jwtManager, "storage", &entity.EmailConfig{
//		SMTPHost:       cfg.SMTP.Host,
//		SMTPPort:       cfg.SMTP.Port,
//		SenderEmail:    cfg.SMTP.Email,
//		SenderPassword: cfg.SMTP.Password,
//	}, txMtx)
//
//	exportSvc := export.NewExportUseCase(authorSvc, bookSvc, "temp")
//	// Инициализация обработчиков
//	bookHandler := v1.NewBookHandler(bookSvc)
//	authorHandler := v1.NewAuthorHandler(authorSvc)
//	authHandler := v1.NewAuthHandler(userSvc)
//	exelHandler := v1.NewExportHandler(exportSvc)
//	// Настройка маршрутов
//	router := gin.Default()
//	routes.SetUpRoutes(router, bookHandler, authorHandler, authHandler, exelHandler, jwtManager)
//
//	return &App{Router: router}, nil
//}

package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"test_go/internal/di"
	"test_go/pkg/transactional"

	"test_go/pkg/httpserver"
	"test_go/pkg/logger"
	"test_go/pkg/postgres"

	"github.com/gin-gonic/gin"

	"test_go/config"

	"test_go/internal/controller/http"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Transaction builder
	pgTx := transactional.NewPgTransaction(pg)

	// RabbitMQ RPC Client
	rmqClient, err := client.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, cfg.RMQ.ClientExchange, cfg.App.Name, cfg.RMQ.ClientPrefix)
	if err != nil {
		l.Fatal("RabbitMQ RPC Client - init error - client.New")
	}

	// Repo
	repo := di.NewRepo(pg)

	// Use-Case
	uc := di.NewUseCase(pgTx, repo, l, cfg)

	// HTTP Server
	handler := gin.New()
	http.NewRouter(handler, cfg, l, uc)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-rmqServer.Notify():
		l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	err = rmqServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	}

	err = rmqClient.Shutdown()
	if err != nil {
		l.Fatal("RabbitMQ RPC Client - shutdown error - rmqClient.RemoteCall", err)
	}
}

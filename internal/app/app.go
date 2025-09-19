package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"test_go/internal/di"
	"test_go/pkg/jwt"
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

	jwtManager := jwt.NewJWTManager(cfg.JWT.SecretKey)

	// Transaction builder
	pgTx := transactional.NewPgTransaction(pg)

	// Repo
	repo := di.NewRepo(pg)

	// Use-Case
	uc := di.NewUseCase(pgTx, repo, l, cfg, jwtManager)

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
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

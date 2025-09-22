package http

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"test_go/internal/controller/http/middleware"
	v1 "test_go/internal/controller/http/v1"
	"test_go/internal/di"
	"test_go/pkg/jwt"
	"test_go/pkg/logger"

	"test_go/config"
	//_ "github.com/finance/fileService/docs" // Swagger docs.
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @description Using a translation service as an example
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(handler *gin.Engine, cfg *config.Config, l logger.Interface, uc *di.UseCase) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	jwtManager := jwt.NewJWTManager(cfg.JWT.SecretKey)
	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	if cfg.Metrics.Enabled {
		handler.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	// Swagger
	if cfg.Swagger.Enabled {
		handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	//Routers
	publicV1Group := handler.Group("/v1")
	{
		v1.NewAuthRoutes(publicV1Group, l, uc.Auth, jwtManager)
	}

	privateV1Group := handler.Group("/v1")
	privateV1Group.Use(middleware.JwtAuthMiddleware(uc.Auth))
	{
		v1.NewUserRoutes(privateV1Group, l, uc.User, jwtManager)
		v1.NewExportRoutes(privateV1Group, l, uc.Export, jwtManager)
		v1.NewBookRoutes(publicV1Group, l, uc.Book, jwtManager)
		v1.NewAuthorRoutes(publicV1Group, l, uc.Author, jwtManager)
	}
}

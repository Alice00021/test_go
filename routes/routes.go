package routes

import (
	"github.com/gin-gonic/gin"
	"test_go/internal/controller/http/v1"

	"test_go/pkg/jwt"
	"test_go/pkg/middleware"
)

func SetUpRoutes(router *gin.Engine, bookHandler *v1.BookHandler, authorHandler *v1.AuthorHandler, authHandler *v1.AuthHandler, exelHandler *v1.ExportHandler, jwtManager *jwt.JWTManager) {
	authGroup := router.Group("auth")

	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.GET("/verify", authHandler.VerifyEmail)
	}
	api := router.Group("/api")
	{
		api.POST("/books", bookHandler.CreateBook)
		api.PATCH("/books/:id", bookHandler.UpdateBook)
		api.DELETE("/books/:id", bookHandler.DeleteBook)
		api.GET("/books/:id", bookHandler.GetBook)
		api.GET("/books", bookHandler.GetAllBooks)

		api.POST("/authors", authorHandler.CreateAuthor)
		api.PATCH("/authors/:id", authorHandler.UpdateAuthor)
		api.DELETE("/authors/:id", authorHandler.DeleteAuthor)
		api.GET("/authors/:id", authorHandler.GetAuthor)
		api.GET("/authors", authorHandler.GetAllAuthors)
		api.GET("/export/statistics", exelHandler.GenerateExportFile)
	}

	protectedAuth := router.Group("/auth")
	protectedAuth.Use(middleware.AuthMiddleware(jwtManager))
	{
		protectedAuth.GET("/profile/:id", authHandler.GetProfile)
		protectedAuth.PATCH("/change-password", authHandler.ChangePassword)
		protectedAuth.PUT("/photo", authHandler.SetProfilePhoto)
		protectedAuth.PATCH("/:id/rating", authHandler.UpdateRating)
		protectedAuth.PATCH("/:id/concurrent-test", authHandler.SimulateConcurrentUpdates)
	}

}

package routes

import (
	"github.com/gin-gonic/gin"

	"test_go/internal/controller/http"
	"test_go/pkg/jwt"
	"test_go/pkg/middleware"
)

func SetUpRoutes(router *gin.Engine, bookHandler *http.BookHandler, authorHandler *http.AuthorHandler, authHandler *http.AuthHandler, jwtManager *jwt.JWTManager){
	authGroup := router.Group("auth")

	{	
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
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
	}
	
	authGroup.Use(middleware.AuthMiddleware(jwtManager))
{
	authGroup.GET("/profile/:id", authHandler.GetProfile)
	authGroup.PATCH("/profile", authHandler.ChangePassword)
	/* authGroup.POST("/logout", handler.Logout) */
}

}


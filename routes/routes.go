package routes

import (
	"github.com/gin-gonic/gin"
	"test_go/internal/controller/http"
)

func SetUpRoutes(router *gin.Engine, bookHandler *http.BookHandler){
	api := router.Group("/api")
	{
		api.POST("/books", bookHandler.CreateBook)
		api.PUT("/books/:id", bookHandler.UpdateBook)
		api.DELETE("/books/:id", bookHandler.DeleteBook)
		api.GET("/books/:id", bookHandler.GetBook)
		api.GET("/books", bookHandler.GetAllBooks)
	}
}


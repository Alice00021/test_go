package routes

import (
	"github.com/gin-gonic/gin"
	"test_go/internal/controller/http"
)

func SetUpRoutes(router *gin.Engine, bookHandler *http.BookHandler, authorHandler *http.AuthorHandler){
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

}


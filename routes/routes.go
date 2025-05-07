package routes

import (
    "test_go/controllers"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {

    api := router.Group("/api")
    {
        api.POST("/authors", controllers.AuthorsCreate(db))
        api.POST("/books", controllers.BookCreate(db))

        api.PUT("/authors/:id", controllers.AuthorUpdate(db))
        api.PUT("/books/:id", controllers.BookUpdate(db))

        api.DELETE("/authors/:id", controllers.AuthorDelete(db))
        api.DELETE("/books/:id", controllers.BookDelete(db))

        api.GET("/authors", controllers.GetAllAuthors(db))
	    api.GET("/authors/:id", controllers.GetOneAuthor(db))

        api.GET("/books/:id", controllers.GetOneBook(db))
        api.GET("/books", controllers.GetAllBooks(db))
       
    }
}


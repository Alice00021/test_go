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
    }
}
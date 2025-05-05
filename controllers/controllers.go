package controllers

import (
    "net/http"

    "test_go/models"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func AuthorsCreate(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var body struct {
            Name   string `json:"name" binding:"required"`
            Female bool   `json:"female"`
        }

        if err := c.ShouldBindJSON(&body); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных: " + err.Error()})
            return
        }

        author := models.Author{Name: body.Name, Female: body.Female}
        result := db.Create(&author)

        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка создания автора: " + result.Error.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"author": author})
    }
}

/* func BookCreate(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var body struct {
			Title      string   `json:"title" binding:"required"`
        }

        if err := c.ShouldBindJSON(&body); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных: " + err.Error()})
            return
        }

        book := models.Author{Name: body.Title}
        result := db.Create(&book)

        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка создания автора: " + result.Error.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"author": book})
    }
} */
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
            Gender bool   `json:"gender"`
        }

        if err := c.ShouldBindJSON(&body); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных: " + err.Error()})
            return
        }

        author := models.Author{Name: body.Name, Gender: body.Gender}
        result := db.Create(&author)

        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка создания автора: " + result.Error.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"author": author})
    }
}

func BookCreate(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var body struct {
			Title      string   `json:"title" binding:"required"`
            AuthorID   uint     `json:"authorID" binding:"required"`
        }

        if err := c.ShouldBindJSON(&body); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных: " + err.Error()})
            return
        }

        book := models.Book{Title:body.Title, AuthorID: body.AuthorID}
        result := db.Create(&book)

        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка создания книги, "  +result.Error.Error()})
            return
        }
        if err :=db.Preload("Author").First(&book, book.ID).Error; err!=nil{
            c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка загрузки данных автора: " + err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"book":book})
    }
}


func AuthorUpdate(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        
        var body struct {
            Name   string `json:"name" binding:"required"`
            Gender bool   `json:"gender"`
        }

        if err := c.ShouldBindJSON(&body); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных: " + err.Error()})
            return
        }
        var author models.Author
        if err := db.Where("id = ?", c.Param("id")).First(&author).Error; err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Запись не найдена!"})
            return
        }

        author.Name = body.Name
        author.Gender = body.Gender

       if err := db.Save(&author).Error; err !=nil{
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка обновления автора: " + err.Error()})
            return
       }

        c.JSON(http.StatusOK, gin.H{"update_author": author})
    }
}

func BookUpdate(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        
        var body struct {
            Title  string `json:"title" binding:"required"`
            AuthorID uint  `json:"authorID" binding:"required"`
        }

        if err := c.ShouldBindJSON(&body); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных: " + err.Error()})
            return
        }
        var book models.Book
        err := db.Transaction(func(tx *gorm.DB) error{
            if err := db.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
                return err
            }
            book.Title = body.Title
            book.AuthorID = body.AuthorID
    
           if err := db.Save(&book).Error; err !=nil{
                return err
           }
            return nil
        })

        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка обновления книги: " + err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"update_book": book})
    }
}

func AuthorDelete(db *gorm.DB) gin.HandlerFunc{
    return  func(c * gin.Context){
        
    }

}

func BookDelete(db *gorm.DB) gin.HandlerFunc{
    return  func(c * gin.Context){
        
    }

}

func AuthorRead(db *gorm.DB) gin.HandlerFunc{
    return  func(c * gin.Context){
        
    }

}

func BookRead(db *gorm.DB) gin.HandlerFunc{
    return  func(c * gin.Context){
        
    }

}




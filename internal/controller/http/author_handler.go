package http

import (
	"net/http"
	"strconv"
	"test_go/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthorHandler struct{
	authorService service.AuthorService
}

func NewAuthorHandler(authorService service.AuthorService) *AuthorHandler{
	return &AuthorHandler{authorService: authorService}
}

func (h *AuthorHandler) CreateAuthor(c *gin.Context){

	var body struct {
		Name   string `json:"name" binding:"required"`
		Gender bool   `json:"gender"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных: " + err.Error()})
		return
	}

	author, err := h.authorService.CreateAuthor(c.Request.Context(), body.Name, body.Gender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"author":author})
}

func (h *AuthorHandler) UpdateAuthor(c *gin.Context){

	var body struct {
		Name   string `json:"name" binding:"required"`
		Gender bool   `json:"gender"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных: " + err.Error()})
		return
	}
	author, err := h.authorService.UpdateAuthor(c.Request.Context(), body.Name, body.Gender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка обновления книги: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"update_author": author})
}

func (h *AuthorHandler) DeleteAuthor(c *gin.Context){
	id_in_type_Str := c.Param("id")
	id, err := strconv.ParseUint(id_in_type_Str, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id пустое"})
		return
	}

        if err := h.authorService.DeleteAuthor(c.Request.Context(), uint(id)); err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

        c.JSON(http.StatusOK, gin.H{
            "message": "Автор был удален успешно",
            "id":      id,
        })
}

func (h *AuthorHandler) GetAuthor(c *gin.Context){
	id_in_type_Str := c.Param("id")
	id, err := strconv.ParseUint(id_in_type_Str, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id пустое"})
		return
	}
	book, err := h.authorService.GetByIDAuthor(c.Request.Context(), uint(id))
	if err !=nil{
		c.JSON(http.StatusNotFound, gin.H{"error":"Автор не найден"})
	}
		c.JSON(http.StatusOK, gin.H{"author": book})
}

func (h *AuthorHandler) GetAllAuthors(c *gin.Context){
	books, err := h.authorService.GetAllAuthors(c.Request.Context())
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"authors":books})
}





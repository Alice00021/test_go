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

// CreateAuthor создает нового автора
// @Summary Создать автора
// @Description Создает нового автора с указанным именем и полом
// @Tags authors
// @Accept json
// @Produce json
// @Param body body object true "Данные автора" { "name": "string", "gender": "boolean" }
// @Success 201 {object} map[string]interface{} "author: созданный автор"
// @Failure 400 {object} map[string]interface{} "неверный формат данных"
// @Failure 500 {object} map[string]interface{} "ошибка сервера"
// @Router /authors [post]
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

// UpdateAuthor обновляет существующего автора
// @Summary Обновить автора
// @Description Обновляет данные автора (имя и пол)
// @Tags authors
// @Accept json
// @Produce json
// @Param body body object true "Обновленные данные автора" { "name": "string", "gender": "boolean" }
// @Success 200 {object} map[string]interface{} "update_author: обновленный автор"
// @Failure 400 {object} map[string]interface{} "неверный формат данных"
// @Failure 500 {object} map[string]interface{} "ошибка сервера"
// @Router /authors/{id} [patch]
func (h *AuthorHandler) UpdateAuthor(c *gin.Context){
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}
	var body struct {
		Name   string `json:"name" binding:"required"`
		Gender bool   `json:"gender"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных: " + err.Error()})
		return
	}
	author, err := h.authorService.UpdateAuthor(c.Request.Context(), body.Name, body.Gender, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка обновления автора: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"update_author": author})
}

// DeleteAuthor удаляет автора по ID
// @Summary Удалить автора
// @Description Удаляет автора по указанному ID
// @Tags authors
// @Produce json
// @Param id path int true "ID автора"
// @Success 200 {object} map[string]interface{} "сообщение об успешном удалении"
// @Failure 400 {object} map[string]interface{} "id пустое"
// @Failure 500 {object} map[string]interface{} "ошибка сервера"
// @Router /authors/{id} [delete]
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

// GetAuthor получает автора по ID
// @Summary Получить автора по ID
// @Description Возвращает автора по указанному ID
// @Tags authors
// @Produce json
// @Param id path int true "ID автора"
// @Success 200 {object} map[string]interface{} "author: найденный автор"
// @Failure 400 {object} map[string]interface{} "id пустое"
// @Failure 404 {object} map[string]interface{} "автор не найден"
// @Router /authors/{id} [get]
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

// GetAllAuthors получает список всех авторов
// @Summary Получить всех авторов
// @Description Возвращает список всех авторов
// @Tags authors
// @Produce json
// @Success 200 {object} map[string]interface{} "authors: список авторов"
// @Failure 500 {object} map[string]interface{} "ошибка сервера"
// @Router /authors [get]
func (h *AuthorHandler) GetAllAuthors(c *gin.Context){
	books, err := h.authorService.GetAllAuthors(c.Request.Context())
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"authors":books})
}








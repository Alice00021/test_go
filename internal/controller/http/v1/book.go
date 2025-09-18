package v1

import (
	"net/http"
	"strconv"
	"test_go/internal/usecase/book"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	bookService book.BookService
}

func NewBookHandler(bookService book.BookService) *BookHandler {
	return &BookHandler{bookService: bookService}
}

// CreateBook создает новую книгу
// @Summary Создать книгу
// @Description Создает новую книгу с указапнным названием и id автора
// @Tags books
// @Accept json
// @Produce json
// @Param body body object true "Данные книги" { "title": "string", "authorID": "uint" }
// @Success 201 {object} map[string]interface{} "book: созданная книга"
// @Failure 400 {object} map[string]interface{} "неверный формат данных"
// @Failure 500 {object} map[string]interface{} "ошибка сервера"
// @Router /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {

	var body struct {
		Title    string `json:"title" binding:"required"`
		AuthorID uint   `json:"authorID" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных: " + err.Error()})
		return
	}
	book, err := h.bookService.CreateBook(c.Request.Context(), body.Title, body.AuthorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"book": book})
}

type UpdateBookRequest struct {
	Title    string `json:"title" binding:"required"`
	AuthorID uint   `json:"authorID"`
}

// UpdateBook обновляет существующую книгу
// @Summary Обновить книгу
// @Description Обновляет данные книги
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "ID книги"
// @Param body body UpdateBookRequest true "Обновленные данные книги"
// @Success 200 {object} map[string]interface{} "update_author: обновленная книга"
// @Failure 400 {object} map[string]interface{} "неверный формат данных"
// @Failure 404 {object} map[string]interface{} "книга не найден"
// @Failure 500 {object} map[string]interface{} "ошибка сервера"
// @Router /books/{id} [patch]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	_, err = h.bookService.GetByIDBook(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Книга не найдена"})
		return
	}

	var body UpdateBookRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных: " + err.Error()})
		return
	}
	book, err := h.bookService.UpdateBook(c.Request.Context(), body.Title, body.AuthorID, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка обновления книги: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"update_book": book})
}

// DeleteBook удаляет книгу по ID
// @Summary Удалить книгу
// @Description Удаляет книгу по указанному ID
// @Tags books
// @Produce json
// @Param id path int true "ID книги"
// @Success 200 {object} map[string]interface{} "сообщение об успешном удалении"
// @Failure 400 {object} map[string]interface{} "id пустое"
// @Failure 500 {object} map[string]interface{} "ошибка сервера"
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {
	id_in_type_Str := c.Param("id")
	id, err := strconv.ParseUint(id_in_type_Str, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id пустое"})
		return
	}

	if err := h.bookService.DeleteBook(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Книга была удалена успешно",
		"id":      id,
	})
}

// GetBook получает книгу по ID
// @Summary Получить книгу по ID
// @Description Возвращает книгу по указанному ID
// @Tags books
// @Produce json
// @Param id path int true "ID книги"
// @Success 200 {object} map[string]interface{} "book: найденная книга"
// @Failure 400 {object} map[string]interface{} "id пустое"
// @Failure 404 {object} map[string]interface{} "книга не найдена"
// @Router /books/{id} [get]
func (h *BookHandler) GetBook(c *gin.Context) {
	id_in_type_Str := c.Param("id")
	id, err := strconv.ParseUint(id_in_type_Str, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id пустое"})
		return
	}
	book, err := h.bookService.GetByIDBook(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Книга не найдена"})
	}
	c.JSON(http.StatusOK, gin.H{"book": book})
}

// GetAllBooks получает список всех книг
// @Summary Получить все книги
// @Description Возвращает список всех книг
// @Tags books
// @Produce json
// @Success 200 {object} map[string]interface{} "books: список книг"
// @Failure 500 {object} map[string]interface{} "ошибка сервера"
// @Router /books [get]
func (h *BookHandler) GetAllBooks(c *gin.Context) {
	books, err := h.bookService.GetAllBooks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books": books})
}

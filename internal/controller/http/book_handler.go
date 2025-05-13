package http

import (
	"net/http"
	"strconv"
	"test_go/internal/service"

	"github.com/gin-gonic/gin"
)

type BookHandler struct{
	bookService service.BookService
}

func NewBookHandler(bookService service.BookService) *BookHandler{
	return &BookHandler{bookService: bookService}
}

func (h *BookHandler) CreateBook(c *gin.Context){

	var body struct {
		Title      string   `json:"title" binding:"required"`
		AuthorID   uint     `json:"authorID" binding:"required"`
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
	
	c.JSON(http.StatusCreated, gin.H{"book":book})
}

func (h *BookHandler) UpdateBook(c *gin.Context){

	var body struct {
		Title  string `json:"title" binding:"required"`
		AuthorID uint  `json:"authorID" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных: " + err.Error()})
		return
	}
	book, err := h.bookService.UpdateBook(c.Request.Context(), body.Title, body.AuthorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка обновления книги: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"update_book": book})
}

func (h *BookHandler) DeleteBook(c *gin.Context){
	id_in_type_Str := c.Param("id")
	id, err := strconv.ParseUint(id_in_type_Str, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id пустое"})
		return
	}

        if err := h.bookService.DeleteBook(c.Request.Context(), uint(id)); err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

        c.JSON(http.StatusOK, gin.H{
            "message": "Книга была удалена успешно",
            "id":      id,
        })
}

func (h *BookHandler) GetBook(c *gin.Context){
	id_in_type_Str := c.Param("id")
	id, err := strconv.ParseUint(id_in_type_Str, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id пустое"})
		return
	}
	book, err := h.bookService.GetByIDBook(c.Request.Context(), uint(id))
	if err !=nil{
		c.JSON(http.StatusNotFound, gin.H{"error":"Книга не найдена"})
	}
		c.JSON(http.StatusOK, gin.H{"book": book})
}

func (h *BookHandler) GetAllBooks(c *gin.Context){
	books, err := h.bookService.GetAllBooks(c.Request.Context())
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books":books})
}

 

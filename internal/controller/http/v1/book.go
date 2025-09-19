package v1

import (
	"net/http"
	"test_go/internal/controller/http/errors"
	"test_go/internal/controller/http/v1/request"
	"test_go/internal/usecase"
	"test_go/internal/utils"
	httpError "test_go/pkg/httpserver"
	"test_go/pkg/jwt"
	"test_go/pkg/logger"
	"test_go/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type bookRoutes struct {
	l  logger.Interface
	uc usecase.Book
}

func NewBookRoutes(privateGroup *gin.RouterGroup, l logger.Interface, uc usecase.Book, jwtManager *jwt.JWTManager) {
	r := &bookRoutes{l, uc}
	{
		h := privateGroup.Group("/book").Use(
			middleware.AuthMiddleware(jwtManager))
		h.POST("/", r.createBook)
		h.PATCH("/:id", r.updateBook)
		h.DELETE("/:id", r.deleteBook)
		h.GET("/:id", r.getBook)
		h.GET("/", r.getAllBooks)
	}
}

func (r *bookRoutes) createBook(c *gin.Context) {
	var req request.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - createBook")
		errors.ErrorResponse(c, httpError.NewBadRequestBodyError(err))
		return
	}
	res, err := r.uc.CreateBook(c.Request.Context(), req.ToEntity())
	if err != nil {
		r.l.Error(err, "http - v1 - createBook")
		errors.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (r *bookRoutes) updateBook(c *gin.Context) {
	id, err := utils.ParsePathParam(utils.ParseParams{Context: c, Key: "id"}, utils.ParseInt64)
	if err != nil {
		r.l.Error(err, "http - v1 - updateBook")
		errors.ErrorResponse(c, httpError.NewBadPathParamsError(err))
		return
	}

	var req request.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - updateBook")
		errors.ErrorResponse(c, httpError.NewBadRequestBodyError(err))
		return
	}
	inp := req.ToEntity()
	inp.ID = id

	if err = r.uc.UpdateBook(c.Request.Context(), inp); err != nil {
		r.l.Error(err, "http - v1 - updateBook")
		errors.ErrorResponse(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (r *bookRoutes) deleteBook(c *gin.Context) {
	id, err := utils.ParsePathParam(utils.ParseParams{Context: c, Key: "id"}, utils.ParseInt64)
	if err != nil {
		r.l.Error(err, "http - v1 - deleteBook")
		errors.ErrorResponse(c, httpError.NewBadPathParamsError(err))
		return
	}

	err = r.uc.DeleteBook(c.Request.Context(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - deleteBook")
		errors.ErrorResponse(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (r *bookRoutes) getBook(c *gin.Context) {
	id, err := utils.ParsePathParam(utils.ParseParams{Context: c, Key: "id"}, utils.ParseInt64)
	if err != nil {
		r.l.Error(err, "http - v1 - getBook")
		errors.ErrorResponse(c, httpError.NewBadPathParamsError(err))
		return
	}

	res, err := r.uc.GetBook(c.Request.Context(), int64(id))
	if err != nil {
		r.l.Error(err, "http - v1 - getBook")
		errors.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (r *bookRoutes) getAllBooks(c *gin.Context) {
	res, err := r.uc.GetBooks(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - getAllBooks")
		errors.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

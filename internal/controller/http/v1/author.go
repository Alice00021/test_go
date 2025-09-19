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

type authorRoutes struct {
	l  logger.Interface
	uc usecase.Author
}

func NewAuthorRoutes(privateGroup *gin.RouterGroup, l logger.Interface, uc usecase.Author, jwtManager *jwt.JWTManager) {
	r := &authorRoutes{l, uc}
	{
		h := privateGroup.Group("/author").Use(
			middleware.AuthMiddleware(jwtManager))
		h.POST("/", r.createAuthor)
		h.PATCH("/:id", r.updateAuthor)
		h.DELETE("/:id", r.deleteAuthor)
		h.GET("/:id", r.getAuthor)
		h.GET("/", r.getAllAuthors)
	}
}

func (r *authorRoutes) createAuthor(c *gin.Context) {
	var req request.CreateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - createAuthor")
		errors.ErrorResponse(c, httpError.NewBadRequestBodyError(err))
		return
	}
	res, err := r.uc.CreateAuthor(c.Request.Context(), req.ToEntity())
	if err != nil {
		r.l.Error(err, "http - v1 - createAuthor")
		errors.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (r *authorRoutes) updateAuthor(c *gin.Context) {
	id, err := utils.ParsePathParam(utils.ParseParams{Context: c, Key: "id"}, utils.ParseInt64)
	if err != nil {
		r.l.Error(err, "http - v1 - updateAuthor")
		errors.ErrorResponse(c, httpError.NewBadPathParamsError(err))
		return
	}

	var req request.UpdateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - updateAuthor")
		errors.ErrorResponse(c, httpError.NewBadRequestBodyError(err))
		return
	}
	inp := req.ToEntity()
	inp.ID = id

	if err = r.uc.UpdateAuthor(c.Request.Context(), inp); err != nil {
		r.l.Error(err, "http - v1 - updateAuthor")
		errors.ErrorResponse(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (r *authorRoutes) deleteAuthor(c *gin.Context) {
	id, err := utils.ParsePathParam(utils.ParseParams{Context: c, Key: "id"}, utils.ParseInt64)
	if err != nil {
		r.l.Error(err, "http - v1 - deleteAuthor")
		errors.ErrorResponse(c, httpError.NewBadPathParamsError(err))
		return
	}

	err = r.uc.DeleteAuthor(c.Request.Context(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - deleteAuthor")
		errors.ErrorResponse(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (r *authorRoutes) getAuthor(c *gin.Context) {
	id, err := utils.ParsePathParam(utils.ParseParams{Context: c, Key: "id"}, utils.ParseInt64)
	if err != nil {
		r.l.Error(err, "http - v1 - getAuthor")
		errors.ErrorResponse(c, httpError.NewBadPathParamsError(err))
		return
	}

	res, err := r.uc.GetAuthor(c.Request.Context(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - getAuthor")
		errors.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (r *authorRoutes) getAllAuthors(c *gin.Context) {
	res, err := r.uc.GetAuthors(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - getAllAuthors")
		errors.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

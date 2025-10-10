package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test_go/internal/controller/http/errors"
	"test_go/internal/controller/http/v1/request"
	"test_go/internal/usecase"
	httpError "test_go/pkg/httpserver"
	"test_go/pkg/logger"
)

type commandRoutes struct {
	l  logger.Interface
	uc usecase.Command
}

func NewCommandRoutes(privateGroup *gin.RouterGroup, l logger.Interface, uc usecase.Command) {
	r := &commandRoutes{l, uc}
	{
		h := privateGroup.Group("/commands")
		h.POST("/upload", r.updateCommands)
	}
}

func (r *commandRoutes) updateCommands(c *gin.Context) {
	var req request.UploadFileRequest
	if err := c.ShouldBind(&req); err != nil {
		r.l.Error(err, "http - v1 - updateCommands")
		errors.ErrorResponse(c, httpError.NewBadRequestBodyError(err))
		return
	}

	if err := r.uc.UpdateCommands(c.Request.Context(), req.File); err != nil {
		r.l.Error(err, "http - v1 - updateCommands")
		errors.ErrorResponse(c, err)
		return
	}

	c.Status(http.StatusOK)
}

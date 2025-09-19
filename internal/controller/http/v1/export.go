package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test_go/internal/controller/http/errors"
	"test_go/internal/usecase"
	"test_go/pkg/jwt"
	"test_go/pkg/logger"
	"test_go/pkg/middleware"
)

type exportRoutes struct {
	l  logger.Interface
	uc usecase.Export
}

func NewExportRoutes(privateGroup *gin.RouterGroup, l logger.Interface, uc usecase.Export, jwtManager *jwt.JWTManager) {
	r := &exportRoutes{l, uc}
	{
		h := privateGroup.Group("/export").Use(
			middleware.AuthMiddleware(jwtManager))
		h.GET("/statistics", r.generateExportFile)
	}
}

func (r *exportRoutes) generateExportFile(c *gin.Context) {
	file, err := r.uc.GenerateExcelFile(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - generateExportFile")
		errors.ErrorResponse(c, err)
		return
	}
	defer file.Close()

	fileName, err := r.uc.SaveToFile(file)
	if err != nil {
		r.l.Error(err, "http - v1 - generateExportFile")
		errors.ErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, fileName)
}

package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test_go/internal/usecase"
)

type ExportHandler struct {
	exportUC usecase.ExportUseCase
}

func NewExportHandler(exportUC usecase.ExportUseCase) *ExportHandler {
	return &ExportHandler{exportUC: exportUC}
}

func (h *ExportHandler) GenerateExportFile(c *gin.Context) {
	file, err := h.exportUC.GenerateExcelFile(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	fileName, err := h.exportUC.SaveToFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"fileName": fileName})
}

package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"net/http"
	"strconv"
	"test_go/internal/usecase"
)

type ExportHandler struct {
	authorService usecase.AuthorService
	bookService   usecase.BookService
}

func NewExportHandler(authorService usecase.AuthorService, bookService usecase.BookService) *ExportHandler {
	return &ExportHandler{authorService, bookService}
}

func (h *ExportHandler) ExportBooksAndAuthorsToExel(c *gin.Context) {
	authors, err := h.authorService.GetAllAuthors(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	books, err := h.bookService.GetAllBooks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bookCount := make(map[uint]int)
	for _, book := range books {
		bookCount[book.AuthorID]++
	}

	f := excelize.NewFile()
	defer func(f *excelize.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	sheetName := "Authors"
	index, _ := f.NewSheet(sheetName)

	headers := []string{"Id", "Имя", "Пол", "Количесвто книг", "Статус"}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Авторы", cell, header)
	}

	for i, author := range authors {
		row := i + 2
		f.SetCellValue(sheetName, "A"+strconv.Itoa(row), author.ID)
		f.SetCellValue(sheetName, "B"+strconv.Itoa(row), author.Name)
		f.SetCellValue(sheetName, "C"+strconv.Itoa(row), map[bool]string{true: "Мужской", false: "Женский"})
		f.SetCellValue(sheetName, "D"+strconv.Itoa(row), bookCount[author.ID])

		status := "Начинающий писатель"

		if bookCount[author.ID] > 5 {
			status = "Профессионал"
		}

		f.SetCellValue(sheetName, "E"+strconv.Itoa(row), status)
	}

	lastRow := len(authors) + 2
	f.SetCellValue(sheetName, "A1"+strconv.Itoa(lastRow), "Всего книг: ")
	f.SetCellFormula(sheetName, "A"+strconv.Itoa(lastRow), "SUM(D2:D"+strconv.Itoa(lastRow-1)+")")

	f.SetActiveSheet(index)

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=authors_export.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

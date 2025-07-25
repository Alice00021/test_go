package usecase

import (
	"context"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type ExportUseCase interface {
	GenerateExcelFile(context.Context) (*excelize.File, error)
	SaveToFile(*excelize.File) (string, error)
}

type exportUseCase struct {
	authorService AuthorService
	bookService   BookService
	exportPath    string
}

func NewExportUseCase(authorService AuthorService, bookService BookService, exportPath string) ExportUseCase {
	if err := os.MkdirAll(exportPath, 0755); err != nil {
		return nil
	}
	return &exportUseCase{
		authorService: authorService,
		bookService:   bookService,
		exportPath:    exportPath,
	}
}

func (uc *exportUseCase) GenerateExcelFile(ctx context.Context) (*excelize.File, error) {
	authors, err := uc.authorService.GetAllAuthors(ctx)
	if err != nil {
		return nil, err
	}

	books, err := uc.bookService.GetAllBooks(ctx)
	if err != nil {
		return nil, err
	}

	bookCount := make(map[uint]int)
	for _, book := range books {
		bookCount[book.AuthorID]++
	}
	f := excelize.NewFile()

	sheetName := "Authors"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	headers := []string{"ID", "Имя", "Пол", "Количество книг", "Статус"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	for i, author := range authors {
		row := i + 2
		f.SetCellValue(sheetName, "A"+strconv.Itoa(row), author.ID)
		f.SetCellValue(sheetName, "B"+strconv.Itoa(row), author.Name)

		gender := "Женский"
		if author.Gender {
			gender = "Мужской"
		}
		f.SetCellValue(sheetName, "C"+strconv.Itoa(row), gender)
		f.SetCellValue(sheetName, "D"+strconv.Itoa(row), bookCount[author.ID])

		status := "Начинающий писатель"
		if bookCount[author.ID] > 5 {
			status = "Профессионал"
		}
		f.SetCellValue(sheetName, "E"+strconv.Itoa(row), status)
	}
	lastRow := len(authors) + 2
	f.SetCellValue(sheetName, "A"+strconv.Itoa(lastRow), "Всего книг: ")
	f.SetCellFormula(sheetName, "D"+strconv.Itoa(lastRow), "SUM(D2:D"+strconv.Itoa(lastRow-1)+")")

	f.SetActiveSheet(index)
	return f, nil
}

func (uc *exportUseCase) SaveToFile(f *excelize.File) (string, error) {
	fileName := "export_" + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
	filePath := filepath.Join(uc.exportPath, fileName)
	if err := f.SaveAs(filePath); err != nil {
		return "", err
	}
	return fileName, nil
}

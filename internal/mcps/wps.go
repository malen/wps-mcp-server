package mcps

import (
	"aoisoft/wps-mcp-server/internal/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

type Wps interface {
	GetSheets() ([]Sheet, error)
	FindSheet(sheetName string) (Sheet, error)
	AddNewSheet(sheetName string) error
	Save() error
}

type Sheet interface {
	Name() (string, error)
	SetValue(cell string, value any) error
	GetValue(cell string) (string, error)
}

type ExcelizeWps struct {
	file *excelize.File
}

type ExcelizeSheet struct {
	file *excelize.File
	name string
}

func NewExcelizeWps(file *excelize.File) Wps {
	return &ExcelizeWps{file: file}
}

func (wps *ExcelizeWps) GetSheets() ([]Sheet, error) {
	sheetList := wps.file.GetSheetList()
	sheets := make([]Sheet, len(sheetList))
	for i, name := range sheetList {
		sheets[i] = &ExcelizeSheet{file: wps.file, name: name}
	}
	return sheets, nil
}

func (wps *ExcelizeWps) FindSheet(sheetName string) (Sheet, error) {
	index, err := wps.file.GetSheetIndex(sheetName)
	if err != nil {
		return nil, err
	}
	if index < 0 {
		return nil, err
	}
	return &ExcelizeSheet{file: wps.file, name: sheetName}, nil
}

func (wps *ExcelizeWps) AddNewSheet(sheetName string) error {
	_, err := wps.file.NewSheet(sheetName)
	if err != nil {
		return err
	}
	return nil
}

func (wps *ExcelizeWps) Save() error {
	file, err := os.OpenFile(filepath.Clean(wps.file.Path), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	return wps.file.Write(file)
}

func (sheet *ExcelizeSheet) Name() (string, error) {
	return sheet.name, nil
}

func (sheet *ExcelizeSheet) SetValue(cell string, value any) error {
	if err := sheet.file.SetCellValue(sheet.name, cell, value); err != nil {
		return err
	}

	if err := sheet.updateDimension(cell); err != nil {
		return err
	}
	return nil
}

func (sheet *ExcelizeSheet) GetValue(cell string) (string, error) {
	value, err := sheet.file.GetCellValue(sheet.name, cell)
	if err != nil {
		return "", err
	}
	if value == "" {
		formula, err := sheet.file.GetCellFormula(sheet.name, cell)
		if err != nil {
			return "", err
		}
		if formula != "" {
			return sheet.file.CalcCellValue(sheet.name, cell)
		}
	}
	return value, nil
}

// 更新工作表的有效区域
func (sheet *ExcelizeSheet) updateDimension(cell string) error {
	dimension, err := sheet.file.GetSheetDimension(sheet.name)
	if err != nil {
		return err
	}
	col1, row1, col2, row2, err := utils.ParseRange(dimension)
	if err != nil {
		return err
	}

	col3, row3, err := excelize.CellNameToCoordinates(cell)
	if err != nil {
		return err
	}

	if col1 > col3 {
		col1 = col3
	}
	if col2 > col3 {
		col2 = row3
	}

	if row1 > row3 {
		row1 = row3
	}
	if row2 < row3 {
		row2 = row3
	}

	startRange, err := excelize.CoordinatesToCellName(col1, row1)
	if err != nil {
		return err
	}

	endRange, err := excelize.CoordinatesToCellName(col2, row2)
	if err != nil {
		return err
	}

	newDimension := fmt.Sprintf("%s:%s", startRange, endRange)
	return sheet.file.SetSheetDimension(sheet.name, newDimension)
}

func OpenFile(absoluteFilePath string) (Wps, func(), error) {
	book, err := excelize.OpenFile(absoluteFilePath)
	if err != nil {
		return nil, func() {}, err
	}
	wps := NewExcelizeWps(book)
	return wps, func() {
		book.Close()
	}, nil
}

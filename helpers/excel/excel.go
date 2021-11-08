package excel

import (
	"errors"
	"github.com/tealeg/xlsx"
)

func WriteToExcel(data [][]string, filename string) (err error) {
	if err != nil {
		return
	}

	f := xlsx.NewFile()
	sheet, err := f.AddSheet("sheet1")
	for _, row := range data {
		r := sheet.AddRow()
		for _, v := range row {
			r.AddCell().Value = v
		}
	}

	return f.Save(filename)
}

func ReadExcel(filename string) (data [][]string, err error) {
	raws, err := xlsx.FileToSlice(filename)
	if err != nil {
		return
	}

	if len(raws) == 0 {
		err = errors.New("no sheet")
	}

	return raws[0], nil
}

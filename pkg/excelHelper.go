package pkg

import "github.com/xuri/excelize/v2"

func GetCell(col, row int) string {
	cell, _ := excelize.CoordinatesToCellName(col, row)
	return cell
}

func GetColumnName(col int) string {
	name, _ := excelize.ColumnNumberToName(col)
	return name
}

func GetHeaderStyle(f *excelize.File) int {
	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"DDEBF7"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	return style
}

func SetRowStyle(f *excelize.File, sheet string, row, cols int, color string) {
	style, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{color}, Pattern: 1},
	})

	for col := 1; col <= cols; col++ {
		cell := GetCell(col, row)
		f.SetCellStyle(sheet, cell, cell, style)
	}
}

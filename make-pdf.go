package main

import (
	"fmt"
	"sort"

	"github.com/jung-kurt/gofpdf"
)

const (
	pageWidth  = 215.9 // letter 8.5 in == 215.9 mm
	numColumns = 3
	margin     = 10
	gutter     = 4
	colWidth   = (pageWidth - 2*margin - (numColumns-1)*gutter) / numColumns
)

var (
	currentColumn int
	pdf           *gofpdf.Fpdf
	yPosition     float64
)

func makePdf() {
	studentList = make([]*student, len(studentMap))
	ii := 0
	for _, stud := range studentMap {
		studentList[ii] = stud
		ii++
	}
	sort.Sort(studentsByName(studentList))

	pdf = gofpdf.New("P", "mm", "Letter", "")

	pdf.SetAcceptPageBreakFunc(func() bool {
		if currentColumn < numColumns-1 {
			setCol(currentColumn + 1)
			pdf.SetY(yPosition)
			// Start new column, not new page
			return false
		}
		setCol(0)
		return true
	})

	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d/{nb}", pdf.PageNo()),
			"", 0, "C", false, 0, "")
	})

	pdf.SetHeaderFunc(func() {
		// Save ordinate
		yPosition = pdf.GetY()
	})

	pdf.SetMargins(margin, margin, margin)
	pdf.AddPage()
	letterHeader := ""
	for ccc, student := range studentList {
		firstLetter := student.name.last[0:1]
		if ccc == 0 {
			fmt.Printf("FIRST ROW FIRST LETTER: %s, LETTERHEADER: %s\n", firstLetter, letterHeader)
		}
		if letterHeader != firstLetter {
			letterHeader = firstLetter
			fmt.Printf("EVERY FIRST LETTER: %s, LETTERHEADER: %s\n", firstLetter, letterHeader)
			pdf.SetTextColor(255, 255, 255)
			pdf.CellFormat(colWidth, 5, letterHeader+"x", "1", 1, "CM", true, 0, "")
		}
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Arial", "B", 12)
		pdf.CellFormat(colWidth, 5, student.name.last+", "+student.name.first, "1", 1, "LM", false, 0, "")
		for _, parent := range student.parents {
			pdf.SetFont("Arial", "", 10)
			parentTxt := fmt.Sprintf("%s: %s, %s", parent.parentType, parent.name.last, parent.name.first)
			pdf.SetCellMargin(5)
			pdf.CellFormat(colWidth, 5, parentTxt, "1", 1, "LM", false, 0, "")
		}
		pdf.SetCellMargin(1)
	}
	pdf.OutputFileAndClose("ChabotDirectory.pdf")

}

func setCol(col int) {
	currentColumn = col
	x := margin + float64(col)*(colWidth+gutter)
	pdf.SetLeftMargin(x)
	pdf.SetX(x)
}

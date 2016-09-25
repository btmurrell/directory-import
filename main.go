package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	s "strings"
)

func main() {
	// Load a TXT file.
	cwd, err := os.Getwd()
	fmt.Println(cwd)
	dir := "src/github.com/btmurrell/directory-import"
	file, err := os.Open(dir + "/inputData/2016-17.csv")
	check(err)

	// Create a new reader.
	reader := csv.NewReader(bufio.NewReader(file))

	writer := csv.NewWriter(os.Stdout)
	i := 0
	for {
		row, err := reader.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}
		fmt.Printf("iteration #%v\n", i)
		if i == 0 {
			i++
			continue
		}

		record := makeRecord(row)

		// if err := writer.Write(record); err != nil {
		// 	log.Fatalln("error writing record to csv:", err)
		// }

		// Display row.
		// ... Display row length.
		// ... Display all individual elements of the slice.
		fmt.Println(row)
		fmt.Println(record)
		fmt.Printf("# columns: %v\n", len(row))
		for value := range row {
			fmt.Printf("  %v\n", row[value])
		}
		i++
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Fatal(err)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func makeRecord(row []string) *Record {
	parentName := row[10]
	nameSplitIdx := s.Index(parentName, " ")
	var parentFName string
	var parentLName string
	if nameSplitIdx >= 0 {
		parentFName = parentName[0:nameSplitIdx]
		parentLName = parentName[nameSplitIdx:len(parentName)]
	} else {
		parentFName = parentName
		parentLName = ""
	}
	stuName := s.Split(row[0], ", ")
	stuFName := stuName[1]
	stuLName := stuName[0]
	record := &Record{
		FirstName:    parentFName,
		LastName:     parentLName,
		Email:        row[14],
		Room:         row[2],
		Grade:        row[7],
		StuFirstName: stuFName,
		StuLastName:  stuLName,
	}
	return record
}

type Record struct {
	FirstName    string
	LastName     string
	Email        string
	Room         string
	Grade        string
	StuFirstName string
	StuLastName  string
}

package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	s "strings"
	log "github.com/Sirupsen/logrus"
	"fmt"
)

func main() {
	log.SetLevel(log.InfoLevel)

	cwd, err := os.Getwd()
	log.Debug(cwd)

	dir := "."
	file, err := os.Open(dir + "/inputData/2016-17.csv")
	check(err)

	reader := csv.NewReader(bufio.NewReader(file))

	writer := csv.NewWriter(os.Stdout)
	i := 0
	roomMap := make(RoomMap)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		log.Debug("iteration #%v\n", i)
		if i == 0 {
			i++
			continue
		}

		record := makeRecord(row)
		roomMap.Add(record.room, record)
		logRow(row, record)
		i++
	}

	// if err := writer.Write(record); err != nil {
	// 	log.Fatalln("error writing record to csv:", err)
	// }

	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Fatal(err)
	}
	room6, _ := roomMap.Peek("6")
	log.WithFields(log.Fields{
		"length": len(room6),
		"list": room6,
	}).Info("ROOM 6: ")
}

func logRow(row []string, record Record) {
	// Display row.
	// ... Display row length.
	// ... Display all individual elements of the slice.
	log.WithFields(log.Fields{
		"row": row,
	}).Debug("ROW")
	log.WithFields(log.Fields{
		"record": fmt.Sprintf("%+v", record),
	}).Info("RECORD")
	log.Debugf("# columns: %v\n", len(row))
	i := 0
	fieldRow := make(log.Fields, len(row))
	for value := range row {
		fieldRow["f"+fmt.Sprintf("%02d", i)] = row[value]
		i++
	}
	log.WithFields(fieldRow).Info("Row fields")
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func makeRecord(row []string) Record {
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
	record := Record{
		firstName:    parentFName,
		lastName:     parentLName,
		email:        row[14],
		room:         row[2],
		grade:        row[7],
		stuFirstName: stuFName,
		stuLastName:  stuLName,
	}
	return record
}

type Record struct {
	firstName    string
	lastName     string
	email        string
	room         string
	grade        string
	stuFirstName string
	stuLastName  string
}

type RoomMap map[string][]Record

func (r RoomMap) Add(key string, value Record) {
	_, ok := r[key]
	if !ok {
		r[key] = make([]Record, 0, 20)
	}
	r[key] = append(r[key], value)
}

func (r RoomMap) Peek(key string) ([]Record, bool) {
	slice, ok := r[key]
	if !ok || len(slice) == 0 {
		return make([]Record,0), false
	}
	return r[key], true
}


package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io"
	"os"
	s "strings"
)

var rowFieldIndices = new(RowFieldIndices)
var noEmailCount = 0

func main() {

	setup()

	dir := "."
	file, err := os.Open(dir + "/inputData/2016-17.csv")
	check(err)

	reader := csv.NewReader(bufio.NewReader(file))

	roomMap := makeRoomMap(reader)

	log.WithFields(log.Fields{
		"count": noEmailCount,
	}).Info("Rows with no email")

	writeRoomCSVFiles(roomMap)

	room6, _ := roomMap.Peek("6")
	log.WithFields(log.Fields{
		"length": len(room6),
		"list":   room6,
	}).Info("ROOM 6")
}

func setup() {
	rowFieldIndices.studentName = 0
	rowFieldIndices.parentName = 10
	rowFieldIndices.parentEmail = 14
	rowFieldIndices.parentEmailAlt = 15
	rowFieldIndices.room = 2
	rowFieldIndices.grade = 7

	log.SetLevel(log.DebugLevel)
}

func makeRoomMap(reader *csv.Reader) RoomMap {
	i := 0
	roomMap := make(RoomMap)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		log.Debugf("iteration #%v", i)
		if i == 0 {
			i++
			continue
		}

		parentRow := makeParentRow(row)
		logRow(row, parentRow)

		if parentRow[2] == "" {
			noEmailCount++
			log.WithFields(log.Fields{
				"row": row,
			}).Error("DISCARDING ROW: no email found")
		} else {
			// only add parentRows which have email addresses
			roomMap.Add(parentRow[3], parentRow)
		}

		i++
	}

	return roomMap
}

func writeRoomCSVFiles(roomMap RoomMap) {

	for room := range roomMap {

		parents, _ := roomMap.Peek(room)
		grade := parents[0][4]

		dir := "./output/"
		fileName := grade + "-rm" + room + ".csv"
		file, err := os.Create(dir + fileName)
		check(err)

		writer := csv.NewWriter(bufio.NewWriter(file))
		if err := writer.WriteAll(parents); err != nil {
			log.Fatalln("error writing parents to csv:", err)
		}
		writer.Flush()
		if err := writer.Error(); err != nil {
			log.Fatal(err)
		}
	}

}

func logRow(row []string, parentRow []string) {

	log.WithFields(log.Fields{
		"count": len(row),
	}).Debug("# columns")

	log.WithFields(log.Fields{
		"row": row,
	}).Debug("RAW ROW")

	i := 0
	rowFields := make(log.Fields, len(row))
	for value := range row {
		rowFields["f"+fmt.Sprintf("%02d", i)] = row[value]
		i++
	}
	log.WithFields(rowFields).Debug("Row fields")

	log.WithFields(log.Fields{
		"parentRow": parentRow,
	}).Debug("PARENT ROW")
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func makeParentRow(row []string) []string {

	// parentName:
	// This implementation is based on value containing one string of "Firstname Lastname"
	// this splits on the space, takes first part as parentFName and all the rest as parentLName
	// * in case the value has no space, parentFName gets it all, parentLName is blank
	// * in case the value multiple spaces, only the first word goes into parentFName, rest to parentLName
	parentName := row[rowFieldIndices.parentName]
	nameSplitIdx := s.Index(parentName, " ")
	var parentFName string
	var parentLName string
	if nameSplitIdx >= 0 {
		parentFName = parentName[0:nameSplitIdx]
		parentLName = parentName[nameSplitIdx+1:]
	} else {
		parentFName = parentName
		parentLName = ""
	}

	// stuName
	// This implementation is based on value containing one string "Lastname, Firstname"
	// this splits on ", ", breaking out the single field into stuFName and stuLName fields
	stuName := s.Split(row[rowFieldIndices.studentName], ", ")
	stuFName := stuName[1]
	stuLName := stuName[0]

	// grade == 0 is Kindergarten; re-assign where appropriate
	grade := row[rowFieldIndices.grade]
	if grade == "0" {
		grade = "K"
	}

	// if no email, check for alternate
	email := row[rowFieldIndices.parentEmail]
	if email == "" {
		if row[rowFieldIndices.parentEmailAlt] != "" {
			email = row[rowFieldIndices.parentEmailAlt]
			log.WithFields(log.Fields{
				"row": row,
			}).Debug("no primary email found; found alt email")
		}
	}

	var result = make([]string, 7)
	result[0] = parentFName
	result[1] = parentLName
	result[2] = email
	result[3] = row[rowFieldIndices.room]
	result[4] = grade
	result[5] = stuFName
	result[6] = stuLName

	return result
}

type RowFieldIndices struct {
	parentName     int
	studentName    int
	parentEmail    int
	parentEmailAlt int
	grade          int
	room           int
}

type RoomMap map[string][][]string

func (r RoomMap) Add(key string, value []string) {
	_, ok := r[key]
	if !ok {
		r[key] = make([][]string, 0, 20)
	}
	r[key] = append(r[key], value)
}
func (r RoomMap) Peek(key string) ([][]string, bool) {
	slice, ok := r[key]
	if !ok || len(slice) == 0 {
		return make([][]string, 0), false
	}
	return r[key], true
}

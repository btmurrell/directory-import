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

	room1, _ := roomMap.Peek("K-1")
	log.WithFields(log.Fields{
		"length": len(room1),
		"list":   room1,
	}).Info("ROOM 1")
}

func setup() {
	rowFieldIndices.studentName = 0
	rowFieldIndices.parentName = 10
	rowFieldIndices.parentEmail = 14
	rowFieldIndices.parentEmailAlt = 15
	rowFieldIndices.room = 2
	rowFieldIndices.grade = 7

	log.SetLevel(log.InfoLevel)
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
			// key is "grade-room"; this accounts for combo rooms,
			// for example, there's usually a combo grade 4 + grade 5
			// class in the same room.  This will allow creating
			// separate csv files for same room, for each of the grades
			roomMap.Add(parentRow[4]+"-"+parentRow[3], parentRow)
		}

		i++
	}

	return roomMap
}

func writeRoomCSVFiles(roomMap RoomMap) {
	header := []string{"FirstName", "LastName", "email", "room", "grade", "StuFn", "StuLn"}
	dir := "./output/"
	for gradeRoom := range roomMap {

		parents, _ := roomMap.Peek(gradeRoom)
		gradeRoomSplitIdx := s.Index(gradeRoom, "-")
		// key is grade-room, split these out
		grade := gradeRoom[0:gradeRoomSplitIdx]
		room := gradeRoom[gradeRoomSplitIdx+1:]

		fileName := "grade" + grade + "-room" + room + ".csv"
		file, err := os.Create(dir + fileName)
		check(err)
		writer := csv.NewWriter(bufio.NewWriter(file))
		if err := writer.Write(header); err != nil {
			log.Fatalf("error writing header to csv file '%v': %v\n", fileName, err)
		}
		if err := writer.WriteAll(parents); err != nil {
			log.Fatalf("error writing parents to csv file '%v': %v\n", fileName, err)
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
		// this simply makes a field label with 2-digit, 0-padded name
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

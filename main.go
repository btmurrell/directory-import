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

var rowFieldIndices = new(RowFieldIndices)
var noEmailCount = 0

func main() {
	rowFieldIndices.studentName = 0
	rowFieldIndices.parentName = 10
	rowFieldIndices.parentEmail = 14
	rowFieldIndices.parentEmailAlt = 15
	rowFieldIndices.room = 2
	rowFieldIndices.grade = 7

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
		log.Debugf("iteration #%v", i)
		if i == 0 {
			i++
			continue
		}

		record := makeRecord(row)
		logRow(row, record)

		if record.email == "" {
			noEmailCount++
			log.WithFields(log.Fields{
				"row": row,
			}).Error("DISCARDING ROW: no email found")
		} else {
			// only add records which have email addresses
			roomMap.Add(record.room, record)
		}

		i++
	}

	log.WithFields(log.Fields{
		"count": noEmailCount,
	}).Info("Rows with no email")

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
	}).Info("ROOM 6")
}

func logRow(row []string, record Record) {

	log.WithFields(log.Fields{
		"count": len(row),
	}).Debug("# columns")

	log.WithFields(log.Fields{
		"row": row,
	}).Debug("RAW ROW")

	i := 0
	rowFields := make(log.Fields, len(row))
	for value := range row {
		rowFields["f" + fmt.Sprintf("%02d", i)] = row[value]
		i++
	}
	log.WithFields(rowFields).Debug("Row fields")

	log.WithFields(log.Fields{
		"record": record,
	}).Debug("RECORD")

}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func makeRecord(row []string) Record {

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

	record := Record{
		firstName:    parentFName,
		lastName:     parentLName,
		email:        email,
		room:         row[rowFieldIndices.room],
		grade:        grade,
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

type RowFieldIndices struct {
	parentName int
	studentName int
	parentEmail int
	parentEmailAlt int
	grade int
	room int
}

func (r Record) String() string {
	return fmt.Sprintf("%#v", r)
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


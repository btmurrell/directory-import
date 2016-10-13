package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io"
	"os"
	s "strings"
)

var rowFieldIndices = new(fieldIndices)
var noEmailCount = 0
var csvRecordsCount = 0
var processedRecordsCount = 0
var outputDir = "csv-output"

func main() {

	logLevel := flag.String("l", "f", "logging level valid values: p (panic), f (fatal), e (error), w (warn), i (info), d (debug)")
	help := flag.Bool("h", false, "help: Show this message")
	flag.Parse()

	if *help {
		usage()
		os.Exit(0)
	}

	setup(loggerMap[*logLevel])

	var inputFileName string
	if len(flag.Args()) < 1 {
		fmt.Println("\n\tYou must enter a file name to convert.")
		fmt.Println("")
		usage()
		os.Exit(1)
	} else {
		inputFileName = flag.Args()[0]
	}

	rooms := makeRoomMap(&inputFileName)

	writeRoomCSVFiles(rooms)

	fmt.Printf("\nFinished successfully processing %v out of %v rows.\n\n", processedRecordsCount, csvRecordsCount)
	fmt.Println("\nYour file has been converted to multiple csv files for import into my-pta.")
	fmt.Printf("You will find all of the files in a folder named '%v' in this directory.\n\n", outputDir)
}

func setup(logLevel log.Level) {
	log.SetLevel(logLevel)

	rowFieldIndices.studentName = 0
	rowFieldIndices.parentName = 10
	rowFieldIndices.parentEmail = 14
	rowFieldIndices.parentEmailAlt = 15
	rowFieldIndices.room = 2
	rowFieldIndices.grade = 7

	os.Mkdir(outputDir, 0755)
}

func makeRoomMap(inputFileName *string) roomMap {
	dir := "./"
	file, err := os.Open(dir + *inputFileName)
	check(err)

	reader := csv.NewReader(bufio.NewReader(file))

	rooms := make(roomMap)

	for {
		row, errRead := reader.Read()
		if errRead == io.EOF {
			logFin()
			break
		}
		log.Debugf("iteration #%v", csvRecordsCount)
		if csvRecordsCount == 0 {
			// header
			csvRecordsCount++
			continue
		}

		parentRow, errParent := makeParentRow(row)
		if errParent == nil {
			// only add parentRows which have email addresses
			// key is "grade-room"; this accounts for combo rooms,
			// for example, there's usually a combo grade 4 + grade 5
			// class in the same room.  This will allow creating
			// separate csv files for same room, for each of the grades
			gradeRoomKey := parentRow[4] + "-" + parentRow[3]
			rooms.Add(gradeRoomKey, parentRow)
			processedRecordsCount++
		} else {
			discardRow(errParent, row)
		}
		logRow(row, parentRow)
		csvRecordsCount++
	}

	return rooms
}

func writeRoomCSVFiles(rooms roomMap) {
	header := []string{"FirstName", "LastName", "email", "room", "grade", "StuFn", "StuLn"}
	for gradeRoom, parents := range rooms {

		gradeRoomSplitIdx := s.Index(gradeRoom, "-")
		// key is grade-room, split these out
		grade := gradeRoom[0:gradeRoomSplitIdx]
		room := gradeRoom[gradeRoomSplitIdx+1:]

		fileName := "grade" + grade + "-room" + room + ".csv"
		file, err := os.Create(outputDir + "/" + fileName)
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

func logFin() {
	log.WithFields(log.Fields{
		"csvRecords":       csvRecordsCount,
		"processedRecords": processedRecordsCount,
	}).Infof("Finished successfully processing %v out of %v rows.\n", processedRecordsCount, csvRecordsCount)
	log.WithFields(log.Fields{
		"count": noEmailCount,
	}).Info("Rows with no email")
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func discardRow(err error, row []string) {
	if rie, ok := err.(*recordImportError); ok {
		if rie.cause == "NO_EMAIL" {
			noEmailCount++
			log.WithFields(log.Fields{
				"row": row,
			}).Error("DISCARDING ROW: no email found")
		}
	}
}

func makeParentRow(row []string) ([]string, error) {

	email, emailErr := resolveEmail(row)
	if _, ok := emailErr.(*recordImportError); ok {
		return nil, emailErr
	}

	parent := resolveParentName(row)
	student := resolveStudentName(row)
	room := row[rowFieldIndices.room]
	// grade == 0 is Kindergarten; re-assign where appropriate
	grade := row[rowFieldIndices.grade]
	if grade == "0" {
		grade = "K"
	}

	var result = make([]string, 7)
	result[0] = parent.first
	result[1] = parent.last
	result[2] = email
	result[3] = room
	result[4] = grade
	result[5] = student.first
	result[6] = student.last

	return result, nil
}

func resolveParentName(row []string) name {
	// parentName:
	// This implementation is based on value containing one string of "Firstname Lastname"
	// this splits on the space, takes first part as parentFName and all the rest as parentLName
	// * in case the value has no space, parentFName gets it all, parentLName is blank
	// * in case the value multiple spaces, only the first word goes into parentFName, rest to parentLName
	parentName := row[rowFieldIndices.parentName]
	student := resolveStudentName(row)
	var parentFName string
	var parentLName string
	if len(parentName) > 0 {
		nameSplitIdx := s.Index(parentName, " ")
		// if we got a split, break it into fname,lname
		if nameSplitIdx > 0 {
			parentFName = parentName[0:nameSplitIdx]
			parentLName = parentName[nameSplitIdx+1:]
		} else {
			// no split, use the single value of
			parentFName = parentName
			parentLName = "[parent unspecified] student: " + student.last
			log.WithFields(log.Fields{
				"first name": parentFName,
				"last name":  parentLName,
				"row":        row,
			}).Warn("could not identify multi-part parent name, using student's last name instead")
		}

	} else {
		// parentName field empty. use student's name as a fallback
		parentFName = "[parent unspecified] student: " + student.first
		parentLName = "[parent unspecified] student: " + student.last
		log.WithFields(log.Fields{
			"first name": parentFName,
			"last name":  parentLName,
			"row":        row,
		}).Warn("No parent name provided, using student's instead")
	}

	return name{parentFName, parentLName}
}

func resolveStudentName(row []string) name {
	// stuName
	// This implementation is based on value containing one string "Lastname, Firstname"
	// this splits on ", ", breaking out the single field into stuFName and stuLName fields
	stuName := s.Split(row[rowFieldIndices.studentName], ", ")
	stuFName := stuName[1]
	stuLName := stuName[0]
	return name{stuFName, stuLName}
}

func resolveEmail(row []string) (string, error) {
	// if no email, check for alternate
	email := row[rowFieldIndices.parentEmail]
	if len(email) == 0 {
		if len(row[rowFieldIndices.parentEmailAlt]) > 1 {
			email = row[rowFieldIndices.parentEmailAlt]
			log.WithFields(log.Fields{
				"altEmail": email,
				"row":      row,
			}).Warn("no primary email found; using alt email")
		}
	}
	if len(email) == 0 {
		return "", &recordImportError{"NO_EMAIL", "No email address found"}
	}
	return email, nil
}

type fieldIndices struct {
	parentName     int
	studentName    int
	parentEmail    int
	parentEmailAlt int
	grade          int
	room           int
}

type name struct {
	first string
	last  string
}

type roomMap map[string][][]string

func (r roomMap) Add(key string, value []string) {
	_, ok := r[key]
	if !ok {
		r[key] = make([][]string, 0, 20)
	}
	r[key] = append(r[key], value)
}
func (r roomMap) Peek(key string) ([][]string, bool) {
	slice, ok := r[key]
	if !ok || len(slice) == 0 {
		return make([][]string, 0), false
	}
	return r[key], true
}

type recordImportError struct {
	cause string
	msg   string
}

func (e *recordImportError) Error() string {
	return fmt.Sprintf("%d - %s", e.cause, e.msg)
}

var loggerMap = map[string]log.Level{
	"p": log.PanicLevel,
	"f": log.FatalLevel,
	"e": log.ErrorLevel,
	"w": log.WarnLevel,
	"i": log.InfoLevel,
	"d": log.DebugLevel,
}

func usage() {
	fmt.Fprintf(os.Stderr, "\nUsage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\n\t%s filename.csv", os.Args[0])
	fmt.Fprintln(os.Stderr, "\nwhere filename.csv is the input file")
	fmt.Fprintln(os.Stderr, "\noptionally, you may specify these flags")
	fmt.Println("")
	flag.PrintDefaults()
}

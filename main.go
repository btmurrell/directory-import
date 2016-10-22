package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	s "strings"

	log "github.com/Sirupsen/logrus"
)

var rowFieldIndices = new(fieldIndices)
var noEmailCount = 0
var csvRecordsCount = 0
var processedRecordsCount = 0
var outputDir = "csv-output"
var students = make(map[string]*student)

func main() {

	logLevel := flag.String("l", "f", "logging level valid values: p (panic), f (fatal), e (error), w (warn), i (info), d (debug)")
	help := flag.Bool("h", false, "help: Show this message")
	flag.Parse()

	if *help {
		usage(0)
	}

	setup(*logLevel)

	var inputFileName string
	if len(flag.Args()) < 1 {
		fmt.Println("\n\tYou must enter a file name to convert.")
		fmt.Println("")
		usage(1)
	} else {
		inputFileName = flag.Args()[0]
	}

	rooms := makeRoomMap(&inputFileName)

	for k, stu := range students {
		fmt.Printf("STUX: %s-> %s\n", k, stu)
		fmt.Printf("\t parent len: %d\n", len(stu.parents))
		for _, par := range stu.parents {
			fmt.Printf("\t PARX: %s\n", par.String())
		}
	}

	writeRoomCSVFiles(rooms)

	fmt.Printf("\nFinished successfully processing %v out of %v rows.\n\n", processedRecordsCount, csvRecordsCount)
	fmt.Println("\nYour file has been converted to multiple csv files for import into my-pta.")
	fmt.Printf("You will find all of the files in a folder named '%v' in this directory.\n\n", outputDir)

	fmt.Printf("Total number of students: %d\n", len(students))
}

func setup(logLevel string) {

	log.SetLevel(loggerMap[logLevel])

	rowFieldIndices.studentName = 0
	rowFieldIndices.teacher = 1
	rowFieldIndices.room = 2
	rowFieldIndices.primaryPhone = 3
	rowFieldIndices.streetAddress = 4
	rowFieldIndices.city = 5
	rowFieldIndices.zip = 6
	rowFieldIndices.grade = 7
	rowFieldIndices.parentType = 9
	rowFieldIndices.parentName = 10
	rowFieldIndices.parentEmail = 14
	rowFieldIndices.parentEmailAlt = 15

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
		if rie.cause == noEmail {
			noEmailCount++
			log.WithFields(log.Fields{
				"row": row,
			}).Errorf("DISCARDING ROW: %s\n", rie.msg)
		}
	}
}

func makeParentRow(row []string) ([]string, error) {

	par := resolveParent(row)
	stuPtr := resolveStudent(row)
	stuPtr.parents = append(stuPtr.parents, par)

	parent := resolveParentName(row)
	student := resolveStudentName(row)

	email, emailErr := resolveEmail(row)
	if _, ok := emailErr.(*recordImportError); ok {
		return nil, emailErr
	}

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

func msgFromImportError(err error) (int, string) {
	if rie, ok := err.(*recordImportError); ok {
		return rie.cause, rie.msg
	}
	return -1, ""
}

func usage(exitCode int) {
	fmt.Fprintf(os.Stderr, "\nUsage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\n\t%s filename.csv", os.Args[0])
	fmt.Fprintln(os.Stderr, "\nwhere filename.csv is the input file")
	fmt.Fprintln(os.Stderr, "\noptionally, you may specify these flags")
	fmt.Println("")
	flag.PrintDefaults()
	os.Exit(exitCode)
}

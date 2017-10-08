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

var (
	rowFieldIndices       = new(fieldIndices)
	noEmailCount          = 0
	csvFileIndex          = 0
	recordsForImportCount = 0
	outputDir             = "csv-output"
	studentMap            = make(map[string]*student)
	studentList           []*student
	srcVersion            = "2017.10.01"
)

func main() {

	logLevel := flag.String("l", "f", "Logging level: valid values: p (panic), f (fatal), e (error), w (warn), i (info), d (debug)")
	help := flag.Bool("h", false, "Help: Show this message")
	version := flag.Bool("v", false, "Version: Show the version of this program")
	flag.Parse()

	if *help {
		usage(0)
	}

	if *version {
		fmt.Fprintln(os.Stdout, "\n\tVersion", srcVersion)
		os.Exit(0)
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

	ingestFile(&inputFileName)
	rooms := makeRoomMap()
	writeRoomCSVFiles(rooms)

	makePdf()

	logFin()

	fmt.Printf("\nFinished successfully processing %v out of %v rows for CSV import.\n", recordsForImportCount, (csvFileIndex - 1))
	fmt.Printf("(%d rows did not have e-mail addresses.)\n\n", noEmailCount)
	fmt.Println("Your file has been converted to multiple CSV files for import into my-pta.")
	fmt.Printf("You will find all of the files in a folder named '%v' in this current folder.\n\n", outputDir)

	fmt.Printf("Total number of students: %d\n", len(studentMap))
}

func setup(logLevel string) {

	log.SetLevel(loggerMap[logLevel])

	rowFieldIndices.studentName = 0
	rowFieldIndices.teacher = 1
	rowFieldIndices.room = 2
	rowFieldIndices.primaryPhone = 4
	rowFieldIndices.streetAddress = 5
	rowFieldIndices.city = 6
	rowFieldIndices.zip = 7
	rowFieldIndices.grade = 8
	rowFieldIndices.parentType = 11
	rowFieldIndices.parentName = 12
	rowFieldIndices.parentEmail = 16
	rowFieldIndices.parentEmailAlt = 10

}

func ingestFile(inputFileName *string) {
	dir := "./"
	file, err := os.Open(dir + *inputFileName)
	check(err)

	reader := csv.NewReader(bufio.NewReader(file))

	for {
		row, errRead := reader.Read()
		if errRead == io.EOF {
			break
		}
		if csvFileIndex == 0 {
			// header
			csvFileIndex++
			continue
		}
		processRow(&row)
		csvFileIndex++
	}
}

func processRow(row *[]string) {
	parPtr := resolveParent(row)
	stuPtr := resolveStudent(row)
	stuPtr.parents = append(stuPtr.parents, parPtr)
	logRow(row, stuPtr)
}

func makeRoomMap() *roomMap {
	rooms := make(roomMap)
	for _, student := range studentMap {
		key := student.gradeVal() + "-" + student.room
		for _, parent := range student.parents {
			if parent.hasEmailError() {
				discardParentFromImport(parent, student)
			} else {
				rooms.add(key, []string{parent.name.first, parent.name.last, parent.email, student.room, student.gradeVal(), student.name.first, student.name.last})
				recordsForImportCount++
			}
		}
	}
	return &rooms
}

func writeRoomCSVFiles(rooms *roomMap) {
	os.Mkdir(outputDir, 0755)
	header := []string{"FirstName", "LastName", "email", "room", "grade", "StuFn", "StuLn"}
	for gradeRoom, parents := range *rooms {
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

func logRow(row *[]string, stud *student) {

	log.Debugf("---------- row #%v --------", csvFileIndex)

	log.WithFields(log.Fields{
		"count": len(*row),
	}).Debug("# columns")

	log.WithFields(log.Fields{
		"row": *row,
	}).Debug("RAW ROW:")

	rowFields := make(log.Fields, len(*row))
	i := 0
	for value := range *row {
		// this simply makes a field label with 2-digit, 0-padded name
		rowFields["f"+fmt.Sprintf("%02d", i)] = (*row)[value]
		i++
	}
	log.WithFields(rowFields).Debug("Row fields:")

	log.WithFields(log.Fields{
		"student": (*stud).String(),
	}).Debug("STUDENT + PARENTS:")
}

func logFin() {
	if log.GetLevel() == log.DebugLevel {
		for k, stu := range studentMap {
			fmt.Printf("STUDENT: %s-> %s\n", k, stu)
			fmt.Printf("\t parent len: %d\n", len(stu.parents))
			for _, par := range stu.parents {
				fmt.Printf("\t PARENT: %s\n", par.String())
			}
		}
	}
	log.WithFields(log.Fields{
		"csvRecords":                csvFileIndex,
		"processedRecordsForImport": recordsForImportCount,
	}).Infof("Finished successfully processing %v out of %v rows.\n", recordsForImportCount, csvFileIndex)
	log.WithFields(log.Fields{
		"count": noEmailCount,
	}).Info("Rows with no email")
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func discardParentFromImport(par *parent, stu *student) {
	log.WithFields(log.Fields{
		"parent":  (*par).String(),
		"student": (*stu).String(),
	}).Errorf("DISCARDING PARENT FROM CSV IMPORT; NO EMAIL ADDRESS")
	noEmailCount++
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
	fmt.Fprintln(os.Stderr, "\n\nwhere filename.csv is the input file")
	fmt.Fprintln(os.Stderr, "\noptionally, you may specify these flags")
	fmt.Println("")
	flag.PrintDefaults()
	os.Exit(exitCode)
}

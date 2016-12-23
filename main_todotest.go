package main

import (
	"reflect"
	"testing"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"test log level"},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_setup(t *testing.T) {
	type args struct {
		logLevel string
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup(tt.args.logLevel)
		})
	}
}

func Test_ingestFile(t *testing.T) {
	type args struct {
		inputFileName *string
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ingestFile(tt.args.inputFileName)
		})
	}
}

func Test_processRow(t *testing.T) {
	type args struct {
		row *[]string
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processRow(tt.args.row)
		})
	}
}

func Test_makeRoomMap(t *testing.T) {
	tests := []struct {
		name string
		want *roomMap
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeRoomMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeRoomMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeFamilyMap(t *testing.T) {
	tests := []struct {
		name string
		want *familyMap
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeFamilyMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeFamilyMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_writeRoomCSVFiles(t *testing.T) {
	type args struct {
		rooms *roomMap
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeRoomCSVFiles(tt.args.rooms)
		})
	}
}

func Test_logRow(t *testing.T) {
	type args struct {
		row  *[]string
		stud *student
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logRow(tt.args.row, tt.args.stud)
		})
	}
}

func Test_logFin(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logFin()
		})
	}
}

func Test_check(t *testing.T) {
	type args struct {
		e error
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check(tt.args.e)
		})
	}
}

func Test_discardParentFromImport(t *testing.T) {
	type args struct {
		par *parent
		stu *student
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discardParentFromImport(tt.args.par, tt.args.stu)
		})
	}
}

func Test_msgFromImportError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := msgFromImportError(tt.args.err)
			if got != tt.want {
				t.Errorf("msgFromImportError() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("msgFromImportError() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_usage(t *testing.T) {
	type args struct {
		exitCode int
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usage(tt.args.exitCode)
		})
	}
}

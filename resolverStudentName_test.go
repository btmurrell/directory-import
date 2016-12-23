package main

import (
	"errors"
	"testing"
)

var nameResp name
var studentNameTests = []struct {
	name     string
	row      *[]string
	expected name
	err      error
}{
	{
		"valid student name",
		&[]string{"Butler, Abigail", "Israel, S-PM", "1", "(510) 827-8437", "4377 MORELAND DR", "CASTRO VALLE", "94546", "0", "Melissa and Robert Butler", "Father", "Robert Butler", "4377 Moreland Dr", "Castro Valley", "94546", "robert@xmzt.org", "melissa@xmzt.org", ""},
		name{"Abigail", "Butler"},
		nil,
	},
	{
		"no comma in student name field",
		&[]string{"Butler Abigail", "Israel, S-PM"},
		nameResp,
		errors.New("no comma found in student name field"),
	},
}

func Test_resolveStudentName(t *testing.T) {
	for _, tt := range studentNameTests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := resolveStudentName(tt.row)
			if err != nil && tt.err == nil {
				t.Errorf("Expected success, got error %v", err)
			} else {
				if actual.first != tt.expected.first {
					t.Errorf("Expected first name of '%s', but was '%s'", tt.expected.first, actual.first)
				}
				if actual.last != tt.expected.last {
					t.Errorf("Expected last name of '%s', but was '%s'", tt.expected.last, actual.last)
				}
			}
		})
	}

}

package main

import (
	"testing"
)

func TestSuccessfulResolveStudentName(t *testing.T) {
	t.Log("Testing successful student name")
	row := []string{"Butler, Abigail", "Israel, S-PM", "1", "(510) 827-8437", "4377 MORELAND DR", "CASTRO VALLE", "94546", "0", "Melissa and Robert Butler", "Father", "Robert Butler", "4377 Moreland Dr", "Castro Valley", "94546", "robert@xmzt.org", "melissa@xmzt.org", ""}
	setup("f")

	expectedName := name{first: "Abigail", last: "Butler"}
	actualName := resolveStudentName(row)

	if actualName.first != expectedName.first {
		t.Errorf("Expected first name of '%s', but was '%s'", expectedName.first, actualName.first)
	}
	if actualName.last != expectedName.last {
		t.Errorf("Expected last name of '%s', but was '%s'", expectedName.last, actualName.last)
	}
}

package main

import (
	log "github.com/Sirupsen/logrus"
	s "strings"
)

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

func resolveParent(row []string) parent {
	pName := resolveParentName(row)
	email, err := resolveEmail(row)
	par := parent{
		name: pName,
		address: address{
			row[rowFieldIndices.streetAddress],
			row[rowFieldIndices.city],
			row[rowFieldIndices.zip],
		},
		primaryPhone: row[rowFieldIndices.primaryPhone],
		parentType:   row[rowFieldIndices.parentType],
	}
	if err != nil {
		if rie, ok := err.(*recordImportError); ok {
			par.meta = make([]*recordImportError, 0)
			par.meta = append(par.meta, rie)
		}
	} else {
		par.email = email
	}
	return par
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

func resolveStudent(row []string) *student {
	stuName := resolveStudentName(row)
	studentCandidate := student{
		name:    stuName,
		teacher: row[rowFieldIndices.teacher],
		room:    row[rowFieldIndices.room],
		grade:   row[rowFieldIndices.grade],
	}
	studentCandidate.parents = make([]parent, 0, 2)

	key := studentCandidate.Key()
	_, ok := students[key]
	if !ok {
		students[key] = &studentCandidate
	}
	return students[key]
}


package main

type family struct {
	studentMap map[string]*student
	parentMap  map[string]*parent
}

func (f *family) String() string {
	resp := "parents: ["
	for _, _parent := range f.parentMap {
		resp = resp + _parent.String() + ", "
	}
	resp = resp + "], students: ["
	for _, _student := range f.studentMap {
		resp = resp + _student.String() + ", "
	}
	resp = resp + "]"
	return resp
}

type familyMap map[string]*family


package main

type family struct {
	studentMap map[string]*student
	parentMap  map[string]*parent
	lastName   string
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

type familiesByName []*family

func (f familiesByName) Len() int      { return len(f) }
func (f familiesByName) Swap(i, j int) { f[i], f[j] = f[j], f[i] }
func (f familiesByName) Less(i, j int) bool {
	return f[i].lastName < f[j].lastName
}

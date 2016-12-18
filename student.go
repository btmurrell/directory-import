package main

import (
	"crypto/md5"
	"encoding/hex"
)

type student struct {
	name    name
	teacher string
	room    string
	grade   string
	parents []*parent
}

func (stu *student) String() string {
	resp := stu.name.String() + ", " + stu.teacher + ", " + stu.room + ", " + stu.grade + ", parents: ["
	for i, par := range stu.parents {
		if i > 0 {
			resp += ", "
		}
		resp += (*par).String()
	}
	resp += "]"
	return resp
}
func (stu *student) uniqueAttributes() string {
	return stu.name.String() + stu.room + stu.grade + stu.teacher
}
func (stu *student) key() string {
	data := []byte(stu.uniqueAttributes())
	sum := md5.Sum(data)
	key := hex.EncodeToString(sum[:md5.Size])
	return key
}
func (stu *student) gradeVal() string {
	if stu.grade == "0" {
		return "K"
	}
	return stu.grade
}

type studentsByName []*student

func (s studentsByName) Len() int      { return len(s) }
func (s studentsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s studentsByName) Less(i, j int) bool {
	iName := s[i].name.last + s[i].name.first
	jName := s[j].name.last + s[j].name.first
	return iName < jName
}

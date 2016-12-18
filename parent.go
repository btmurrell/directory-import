package main

import (
	"crypto/md5"
	"encoding/hex"
)

type parent struct {
	name         name
	address      address
	primaryPhone string
	parentType   string
	studentKeys  []string
	email        string
	meta         []*recordImportError
}

func (par *parent) String() string {
	resp := par.parentType + ": " + par.name.String() + ", " + par.address.String() + ", " + par.email + ", " + par.primaryPhone
	//+ ", students: ["
	//for i, stu := range par.students {
	//	if i > 0 {
	//		resp += ", "
	//	}
	//	resp += (*stu).String()
	//}
	//resp += "]"
	if len(par.meta) > 0 {
		resp += ", meta: ["
		for i, err := range par.meta {
			_, msg := msgFromImportError(err)
			if i > 0 {
				resp += ", "
			}
			resp += "ERROR: " + msg
		}
		resp += "]"
	}
	return resp
}
func (par *parent) uniqueAttributes() string {
	return par.name.String() + par.email
}
func (par *parent) key() string {
	data := []byte(par.uniqueAttributes())
	sum := md5.Sum(data)
	key := hex.EncodeToString(sum[:md5.Size])
	return key
}
func (par *parent) hasEmailError() bool {
	hasError := false
	if len(par.meta) > 0 {
		for _, err := range par.meta {
			errType, _ := msgFromImportError(err)
			if errType == noEmail {
				hasError = true
			}
		}
	}
	return hasError
}

type parentsByName []*parent

func (p parentsByName) Len() int      { return len(p) }
func (p parentsByName) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p parentsByName) Less(i, j int) bool {
	iName := p[i].name.last + p[i].name.first
	jName := p[j].name.last + p[j].name.first
	return iName < jName
}

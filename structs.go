package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

const noEmail = iota

var loggerMap = map[string]log.Level{
	"p": log.PanicLevel,
	"f": log.FatalLevel,
	"e": log.ErrorLevel,
	"w": log.WarnLevel,
	"i": log.InfoLevel,
	"d": log.DebugLevel,
}

type fieldIndices struct {
	parentName     int
	studentName    int
	parentEmail    int
	parentEmailAlt int
	grade          int
	room           int
	teacher        int
	primaryPhone   int
	streetAddress  int
	city           int
	zip            int
	parentType     int
}

type recordImportError struct {
	cause int
	msg   string
}

func (err *recordImportError) Error() string {
	return fmt.Sprintf("[%v] - %s", err.cause, err.msg)
}

type roomMap map[string][][]string

func (r roomMap) Add(key string, value []string) {
	_, ok := r[key]
	if !ok {
		r[key] = make([][]string, 0, 20)
	}
	r[key] = append(r[key], value)
}
func (r roomMap) Peek(key string) ([][]string, bool) {
	slice, ok := r[key]
	if !ok || len(slice) == 0 {
		return make([][]string, 0), false
	}
	return r[key], true
}

type name struct {
	first string
	last  string
}
func (n *name) String() string {
	return n.last + ", " + n.first
}

type student struct {
	name    name
	teacher string
	room    string
	grade   string
	parents []parent
}

func (stu *student) String() string {
	return stu.name.String() + ", " + stu.teacher + ", " + stu.room + ", " + stu.grade
}
func (stu student) Key() string {
	data := []byte(stu.String())
	sum := md5.Sum(data)
	key := hex.EncodeToString(sum[:md5.Size])
	return key
}

type address struct {
	street string
	city   string
	zip    string
}

func (a *address) String() string {
	return a.street + ", " + a.city + ", CA " + a.zip
}

type parent struct {
	name         name
	address      address
	primaryPhone string
	parentType   string
	students     []*student
	email        string
	meta         []*recordImportError
}

func (par *parent) String() string {
	resp := par.parentType + ": " + par.name.String() + ", " + par.address.String() + ", " + par.email + ", " + par.primaryPhone
	if len(par.meta) > 0 {
		for _, err := range par.meta {
			_, msg := msgFromImportError(err)
			resp += "\n\t\tERROR " + msg
		}
	}
	return resp
}
func (par *parent) hasEmailError() bool  {
	if len(par.meta) > 0 {
		for _, err := range par.meta {
			errType, _ := msgFromImportError(err)
			return errType == noEmail
		}
	}
	return false

}
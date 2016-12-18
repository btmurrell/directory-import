package main

import (
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

func (r roomMap) add(key string, value []string) {
	_, ok := r[key]
	if !ok {
		r[key] = make([][]string, 0, 20)
	}
	r[key] = append(r[key], value)
}
func (r roomMap) peek(key string) ([][]string, bool) {
	slice, ok := r[key]
	if !ok || len(slice) == 0 {
		return make([][]string, 0), false
	}
	return r[key], true
}

type family struct {
	studentMap map[string]*student
	parentMap  map[string]*parent
}
type familyMap map[string]*family

type name struct {
	first string
	last  string
}

func (n *name) String() string {
	return n.last + ", " + n.first
}

type address struct {
	street string
	city   string
	zip    string
}

func (a *address) String() string {
	return a.street + ", " + a.city + ", CA " + a.zip
}

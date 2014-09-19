package log

import (
	"log"
)

type Level int

/*
const (
	Off Level = iota
	Trace
	Debug
	Info
	Warn
	Error
	Fatal
	All
)

*/

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Println(v ...interface{}) {
	log.Println(v...)
}

/*
func Fatal(v ...interface{}) {
	log.Fatal(v...)
}
*/

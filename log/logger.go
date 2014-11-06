package log

import (
	"github.com/ms140569/ghost/log/level"
	"log"
)

var SystemLogLevel level.Level

func SetSystemLogLevel(l level.Level) {
	SystemLogLevel = l
}

func SetSystemLogLevelFromString(s string) {
	SystemLogLevel = level.ToLoglevel(s)
}

func LeveledLogger(level level.Level, format string, v ...interface{}) {
	if level >= SystemLogLevel {
		log.Printf(level.String()+" : "+format, v...)
	}
}

func Trace(format string, v ...interface{}) {
	LeveledLogger(level.Trace, format, v...)
}

func Debug(format string, v ...interface{}) {
	LeveledLogger(level.Debug, format, v...)
}

func Info(format string, v ...interface{}) {
	LeveledLogger(level.Info, format, v...)
}

func Warn(format string, v ...interface{}) {
	LeveledLogger(level.Warn, format, v...)
}

func Error(format string, v ...interface{}) {
	LeveledLogger(level.Error, format, v...)
}

func Fatal(format string, v ...interface{}) {
	LeveledLogger(level.Fatal, format, v...)
}

func All(format string, v ...interface{}) {
	LeveledLogger(level.All, format, v...)
}

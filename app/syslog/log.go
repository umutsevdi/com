package syslog

import "log"

type LOG_LEVEL int

var baseLogLevel LOG_LEVEL

func SetLogLevel(level LOG_LEVEL) {
	baseLogLevel = level
}

const (
	TRACE   = LOG_LEVEL(0)
	INFO    = LOG_LEVEL(1)
	WARNING = LOG_LEVEL(2)
	ERROR   = LOG_LEVEL(3)
	FATAL   = LOG_LEVEL(4)
)

func String(level LOG_LEVEL) string {
	switch level {
	case LOG_LEVEL(INFO):
		return "INFO    "
	case LOG_LEVEL(WARNING):
		return "WARNING "
	case LOG_LEVEL(ERROR):
		return "ERROR   "
	case LOG_LEVEL(FATAL):
		return "FATAL   "
	}
	return "TRACE   "
}

func write(level LOG_LEVEL, objects ...interface{}) {
	if int(level) >= int(baseLogLevel) {
		if int(level) == 4 {
			log.Fatal(String(level), objects)
		}
		log.Println(String(level), objects)
	}
}

func Trace(objects ...interface{}) {
	write(TRACE, objects...)
}

func Info(objects ...interface{}) {
	write(INFO, objects...)
}

func Warn(objects ...interface{}) {
	write(WARNING, objects...)
}

func Error(objects ...interface{}) {
	write(ERROR, objects...)
}

func Fatal(objects ...interface{}) {
	write(FATAL, objects...)
}

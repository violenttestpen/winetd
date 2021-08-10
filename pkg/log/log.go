package log

import "log"

// ERROR, WARNING and DEBUG are log levels to control verbosity
const (
	ERROR = iota
	WARNING
	INFO
)

// Log represents a logger object
type Log struct {
	logLevel int
}

// NewLogger returns a new logger set at the specified log level
func NewLogger(logLevel int) Log {
	return Log{logLevel: logLevel}
}

func (l *Log) Fatal(msg ...interface{}) {
	log.Fatal(msg...)
}

func (l *Log) Error(msg ...interface{}) {
	if l.logLevel >= ERROR {
		log.Println(msg...)
	}
}

func (l *Log) Warning(msg ...interface{}) {
	if l.logLevel >= WARNING {
		log.Println(msg...)
	}
}

func (l *Log) Info(msg ...interface{}) {
	if l.logLevel >= INFO {
		log.Println(msg...)
	}
}

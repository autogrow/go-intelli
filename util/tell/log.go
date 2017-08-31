package tell

import (
	"fmt"
	"log"
	"os"

	"bitbucket.org/autogrowsystems/easydose/events"
)

// log levels
const (
	DEBUG = 0
	INFO  = 1
	WARN  = 2
	ERROR = 3
	FATAL = 4
)

// Level is the current log level
var Level = 0

func init() {
	log.SetFlags(log.Ldate | log.Ltime)
}

// Debugf logs a Debugf message
func Debugf(msg string, args ...interface{}) {
	if Level > DEBUG {
		return
	}

	msg = fmt.Sprintf(msg, args...)
	log.Printf("DEBUG: " + msg)
}

// Infof logs a Infof message
func Infof(msg string, args ...interface{}) {
	if Level > INFO {
		return
	}

	msg = fmt.Sprintf(msg, args...)
	log.Printf("INFO: " + msg)
}

// Warnf logs a Warnf message
func Warnf(msg string, args ...interface{}) {
	if Level > WARN {
		return
	}

	msg = fmt.Sprintf(msg, args...)
	log.Printf("WARN: " + msg)
}

// IfErrorf logs an error message if there was an error
func IfErrorf(err error, msg string, args ...interface{}) {
	if err != nil {
		msg = fmt.Sprintf(msg, args...)
		Errorf(msg+": %", err)
	}
}

// Errorf logs a Errorf message
func Errorf(msg string, args ...interface{}) {
	if Level > ERROR {
		return
	}

	msg = fmt.Sprintf(msg, args...)
	log.Printf("ERROR: " + msg)
}

// Fatalf logs a Fatalf message
func Fatalf(msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	events.LogNotify("fatal", msg)
	log.Printf("FATAL: " + msg)
	os.Exit(1)
}

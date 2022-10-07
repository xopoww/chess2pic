package chess2pic

import (
	"log"
	"os"
)

var DEBUG bool = false

func Infof(format string, a ...interface{}) {
	log.Printf(format, a...)
}

func Debugf(format string, a ...interface{}) {
	if !DEBUG {
		return
	}
	Infof(format, a...)
}

func Fatalf(format string, a ...interface{}) {
	Infof(format, a...)
	os.Exit(1)
}

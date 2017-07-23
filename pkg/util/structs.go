package util

import (
	"io/ioutil"
	"log"
	"os"
	"time"
)

type (
	// Logger wraps loggers for stdout, stderr and debug output
	Logger struct {
		Out   *log.Logger
		Err   *log.Logger
		Debug *log.Logger
	}
	// A FileStat describes a local and remote file and can contain an error if the information
	// was not possible to get
	FileStat struct {
		Err     error
		Name    string
		Dir     string
		Path    string
		Size    int64
		ModTime time.Time
		//Checksum string
	}
)

// NewLogger creates a new Logger ready for use
func NewLogger(debug, onlyShowErrors bool) *Logger {
	l := &Logger{
		Out:   log.New(os.Stdout, "", 0),
		Err:   log.New(os.Stderr, "", 0),
		Debug: log.New(os.Stdout, "[DEBUG] ", 0),
	}
	if !debug {
		l.Debug.SetOutput(ioutil.Discard)
	}

	if onlyShowErrors {
		l.Debug.SetOutput(ioutil.Discard)
		l.Out.SetOutput(ioutil.Discard)
	}
	return l
}

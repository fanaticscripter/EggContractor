package main

import (
	log "github.com/sirupsen/logrus"
)

var _errored bool

func logError(args ...interface{}) {
	log.Error(args)
	_errored = true
}

func logErrorf(format string, args ...interface{}) {
	log.Errorf(format, args)
	_errored = true
}

package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func initLogger() {
	now := time.Now()
	year, month, day := now.Date()
	logfilename := fmt.Sprintf("log/info.log.%04d_%02d_%02d_%02d", year, month, day, now.Hour())

	_logfile, err := os.OpenFile(logfilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	exitIfErr(err)

	logfile = _logfile
	logger = log.New(logfile, "", log.LstdFlags|log.Lshortfile)
	logger.Println("wschat server start...")
}

func logger_file_update() {
	now := time.Now()
	year, month, day := now.Date()
	logfilename := fmt.Sprintf("log/info.log.%04d_%02d_%02d_%02d", year, month, day, now.Hour())
	_logfile, err := os.OpenFile(logfilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	exitIfErr(err)

	logger = log.New(_logfile, "", log.LstdFlags|log.Lshortfile)
	__logfile := logfile
	logfile = _logfile
	__logfile.Close()
}

package main

import (
	"time"
)

const (
	DAY_TIME  = 24 * 60 * 60
	HOUR_TIME = 60 * 60
	MIN_TIME  = 60
)


func timer() {

	logupdate_ticker := time.NewTicker(DAY_TIME * time.Second)

	for {
		select {
		case <-logupdate_ticker.C:
			logger_file_update()
		}
	}
}

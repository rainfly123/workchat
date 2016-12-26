package main

import (
//"fmt"
)

func exitIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func stoplive(live *Live) {
	for _, conn := range live.connections {
		conn.ws.Close()
	}
	live.connections = live.connections[:0]
}

func pushmsg(live *Live, msg string) {
	for _, conn := range live.connections {
		select {
		case conn.send <- []byte(msg):
		default:
			close(conn.send)
		}
	}
}

package main

import (
	"fmt"
)

func newMainhub() *Mainhub {
	return &Mainhub{
		register:   make(chan *connection),
		unregister: make(chan *connection),
	}
}

func (h *Mainhub) run() {
	for {
		select {
		case c := <-h.register:
			/* add connections in live */
			live, _ := live_map[c.liveid]
			live.connections = append(live.connections, c)
		case c := <-h.unregister:
			logger.Printf("connection %s closed\n", c.name)

			/* delete connections in live */
			liveid := c.liveid
			live, ok := live_map[liveid]
			if ok == true {
				index := -1
				for i, val := range live.connections {
					if val == c {
						index = i
					}
				}
				if index != -1 {
					live.connections = append(live.connections[:index], live.connections[index+1:]...)
				}

				/* info other users */
				if c.login == true {
					info := fmt.Sprintf("%s 离开直播间", c.name)
					pushmsg(live, info)
				}
			}
			close(c.send)
		}
	}
}

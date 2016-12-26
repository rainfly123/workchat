package main

import (
//	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"math/rand"
//	"time"
)

func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		/* if user is been silented */
		if c.silent == true || c.login == false {
			//c.send <- []byte(fmt.Sprintf("%s:被禁言", c.userid))
			c.send <- []byte(fmt.Sprintf("被禁言"))
			continue
		}

		live, ok := live_map[c.liveid]
		if ok == true && !c.silent {
                        tempmy := fmt.Sprintf("%s:%s:",c.userid, c.name)
                        //message = append([]byte(":"), message...)
			message = append([]byte(tempmy), message...)
			//msgrecord := []byte(fmt.Sprintf("%d %s\n", time.Now().Unix(), string(message)))

			//live.chatrecord = append(live.chatrecord, msgrecord...)

			/* record chat messages into redis */
			/*if len(live.chatrecord) > CHAT_BUFF_MAX {
				chatlen := len(live.chatrecord)
				//tmpchatrecord := live.chatrecord[:chatlen]
				live.chatrecord = live.chatrecord[chatlen:]
				//appendchat(c.liveid, tmpchatrecord)
			}*/

			/* send to other user */
			for _, c_chan := range live.connections {
				select {
				case c_chan.send <- message:
				default:
					close(c_chan.send)
				}
			}
		}
	}
//	c.ws.Close()
}

func (c *connection) writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, 
  CheckOrigin: func(r *http.Request) bool { return true }, WriteBufferSize: 1024}

type wsHandler struct {
	h *Mainhub
}

func (wsh wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.Printf("new user try to connect:%s\n", r.URL.RawQuery)

	querys := r.URL.Query()
	var name string = ""
	var userid string = ""
	var liveid string = ""
	var err error
	if querys["name"] != nil && querys["name"][0] != "" {
		name = querys["name"][0]
	}
	if querys["liveid"] != nil && querys["liveid"][0] != "" {
		liveid = querys["liveid"][0]
	}
	if querys["userid"] != nil && querys["userid"][0] != "" {
		userid = querys["userid"][0]
	}
	if liveid == "" {
                fmt.Println("no chat room")
		return
	}

	live, ok := live_map[liveid]
	if ok != true {
 
		live_map[liveid] = &Live{
			connections: make([]*connection, 0, 10000),
			chatrecord:  make([]byte, 0, 1000),
		}
	}
	live, ok = live_map[liveid]
        

	/* check repeat login */
	for _, conn := range live.connections {
		if userid != "" && conn.userid == userid {
			logger.Printf("user(id:%s) repeat login\n", userid)
			return
		}
	}

	/* check origin */
	//r.Header["Origin"][0] = "http://192.168.1.17:8080"
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Println("fail to upgrade connection")
		return
	}

	logger.Printf("new user connected:%s\n", r.URL.RawQuery)

	login := false
	if name != "" && userid != "" {
		login = true
	}

	/* info other users */
	if login == true {
		//info := fmt.Sprintf("%s:login:%s", userid, name)
		info := fmt.Sprintf("%s 进入直播间", name)
		pushmsg(live, info)

                live.online += (rand.Intn(5) + 1) 
		info = fmt.Sprintf("online_%d", live.online)
		pushmsg(live, info)
	}

	/* add client in connections */
	c := &connection{
		send:   make(chan []byte, 256),
		ws:     ws,
		h:      wsh.h,
		name:   name,
		userid: userid,
		liveid: liveid,
		silent: false,
		login:  login,
	}
	c.h.register <- c

	/* check been silent */
	for _, _userid := range live.silentusers {
		if userid == _userid {
			c.silent = true
		}
	}
        /*
	select {
	//case c.send <- []byte(__userlist):
	case c.send <- []byte("登录成功"):
	default:
		close(c.send)
	}
          */

	defer func() { c.h.unregister <- c }()

	go c.writer()
	c.reader()
}

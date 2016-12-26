package main

import (
	"flag"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
	"log"
	"os"
)

type Live struct {
	connections []*connection `json:"-"`
	silentusers []string      `json:"-"`
	chatrecord  []byte        `json:"-"`
	online      int           `json:"-"`
}

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// The Mainhub.
	h *Mainhub

	// name, channel
	name   string
	userid string
	liveid string
	silent bool
	login  bool // a normal user or an anonymous client
}

type User struct {
	Name   string `json:"name"`
	Userid string `json:"userid"`
}

type Userlist struct {
	Creator  User    `json:"creator"`
	Userlist []*User `json:"userlist"`
}

type Mainhub struct {
	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

const (
	ERR_SUCCESS int32 = iota
	ERR_PARAMETERS
	ERR_NOLIVEID
	ERR_REPEATLOGIN
)

var ERROR_MAP = map[int32]string{
	ERR_SUCCESS:     "success",
	ERR_PARAMETERS:  "parameters wrong",
	ERR_NOLIVEID:    "no such live",
	ERR_REPEATLOGIN: "repeat login",
}

const (
	VERSION       string = "1.0.0"
	MAIN_VERSION  uint8  = 1
	MID_VERSION   uint8  = 0
	LAST_VERSION  uint8  = 0
	CHAT_BUFF_MAX int    = 10000
	ThriftAddr    string = ":9090"
	yq_redishost  string = "192.168.1.17:6379"
)

var (
	addr = flag.String("addr", ":8080", "websocket chat service address")

	live_map     map[string]*Live
	logfile      *os.File
	logger       *log.Logger
	chat_rdspool *redis.Pool
)

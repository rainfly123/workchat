package main

import (
	"flag"
	"io"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)


func silentHandle(w http.ResponseWriter, req *http.Request) {
    liveid := req.FormValue("liveid")
    userid := req.FormValue("userid")
    if len(liveid) < 6 {
        if len(userid) < 3 {
           io.WriteString(w, "parameter error!\n")
           return
        }
    }
    live, ok := live_map[liveid]
    if ok == false {
        io.WriteString(w, "no liveid!\n")
        return
    }

    live_map[liveid].silentusers = append(live_map[liveid].silentusers, userid)
    for _, conn := range live.connections {
		if conn.userid == userid {
			conn.silent = true
		}
    }
    io.WriteString(w, "ok!\n")
}

func main() {
	flag.Parse()

	live_map = make(map[string]*Live)

	initLogger()
	//initRedis()

	/* start timer */
	go timer()

	h := newMainhub()
	go h.run()

        go check_persons()
	/* start unix server */
	//go unixServer()

	http.Handle("/wspage/", http.FileServer(http.Dir("./js")))
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./js"))))
	http.Handle("/ws", wsHandler{h: h})
//	http.Handle("/testapi/", testapiHandler{})
        //file upload
        http.HandleFunc("/image", uploadHandle)
        http.HandleFunc("/silent", silentHandle)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		logger.Println("wschat ListenAndServe:", err)
	}
}

package main

import (
        "time"
)
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

const RUNNING = 1

func updateMysqlPersons(liveid string, persons int) {
    db, err := sql.Open("mysql", "root:123321@/live")
    if err != nil {
        panic(err.Error()) 
    }
    defer db.Close()

    // Prepare statement for inserting data
    stmtIns, err := db.Prepare("update live set persons=? where liveid=?") // ? = placeholder
    if err != nil {
        panic(err.Error()) 
    }
    defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

    _, err = stmtIns.Exec(persons, liveid)
    if err != nil {
         panic(err.Error()) // proper error handling instead of panic in your app
    }
}

func checkMysqlState(liveid string)(int) {
    db, err := sql.Open("mysql", "root:123321@/live")
    if err != nil {
        panic(err.Error()) 
    }
    defer db.Close()
    rows:= db.QueryRow("select state from live where liveid=?", liveid)
    if err != nil {
        panic(err.Error()) 
    }
    var state int
    state = 0
    rows.Scan(&state)
    return state
}




func check_persons(){

    for {
            for k, v := range live_map { 

                liveid := k
                //persons:= len(v.connections)
                persons:= v.online

                if checkMysqlState(liveid) == RUNNING {
                    updateMysqlPersons(liveid, persons)
                }
            }
            time.Sleep(time.Second * 60)
    }
}


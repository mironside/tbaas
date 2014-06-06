package main

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func POST_Messages(w http.ResponseWriter, r *http.Request) {
	var message map[string]interface{}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &message)
	fmt.Println(message)
	db, _ := sqlite3.Open(":memory:")
	db.Exec("insert into messages (timestamp, username, text) values (?, ?, ?)", time.Now(), message["username"], message["text"])
	db.Commit()
}

func GET_Messages(w http.ResponseWriter, r *http.Request) {
	db, _ := sqlite3.Open(":memory:")
	messages := []map[string]interface{}{}
	for s, err := db.Query("select * from messages order by timestamp asc"); err == nil; err = s.Next() {
		row := make(sqlite3.RowMap)
		s.Scan(row)
		fmt.Println(row)
		messages = append(messages, row)
	}
	body, _ := json.Marshal(messages)
	w.Write(body)
}

func MessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		POST_Messages(w, r)
	} else if r.Method == "GET" {
		GET_Messages(w, r)
	}
}

func main() {
	db, _ := sqlite3.Open(":memory:")
	db.Exec("create table messages(timestamp, username, text)")

	fmt.Println("talk to me...on port 8080")
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/messages", MessagesHandler)
	http.ListenAndServe(":8081", nil)
}

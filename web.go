package main

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var db = "test.db"

func POST_Messages(w http.ResponseWriter, r *http.Request) {
	var message map[string]interface{}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &message)
	fmt.Println(message)
	db, _ := sqlite3.Open(db)
	db.Exec("insert into messages (timestamp, username, text) values (?, ?, ?)", time.Now(), message["username"], message["text"])
	db.Commit()
}

func GET_Messages(w http.ResponseWriter, r *http.Request) {
	db, _ := sqlite3.Open(db)
	messages := []map[string]interface{}{}

	r.ParseForm()

	sinceId, err := strconv.ParseInt(r.FormValue("since_id"), 10, 64)
	if err != nil {
		sinceId = 0
	}
	maxId, err := strconv.ParseInt(r.FormValue("max_id"), 10, 64)
	if err != nil {
		maxId = 0
	}
	count, err := strconv.ParseInt(r.FormValue("count"), 10, 64)
	if err != nil {
		count = 0
	}

	fmt.Println(sinceId, maxId, count)

	params := []interface{}{}
	query := "select * from (select * from messages "
	if sinceId > 0 || maxId > 0 {
		query += "where "
	}
	if sinceId > 0 {
		query += "id > ? "
		params = append(params, sinceId)
	}
	if maxId > 0 {
		if sinceId > 0 {
			query += " and "
		}
		query += "id <= ? "
		params = append(params, maxId)
	}
	query += "order by id desc "
	if count > 0 {
		query += "limit ? "
		params = append(params, count)
	}
	query += ") sub order by id asc"

	fmt.Println(query, params)

	for s, err := db.Query(query, params...); err == nil; err = s.Next() {
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
	db, _ := sqlite3.Open(db)
	db.Exec("create table messages(id integer primary key autoincrement, timestamp, username, text)")

	fmt.Println("talk to me...on port 8080")
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/messages", MessagesHandler)
	http.ListenAndServe(":8080", nil)
}

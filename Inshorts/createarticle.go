package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

//POSTHandler to insert an article and store in db
func POSTHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	defer db.Close()
	var a Article
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(a.Title) == 0 || len(a.Subtitle) == 0 || len(a.Content) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "message": "Title/Subtitle/Content Cannot be empty }`))
		return
	}
	a.Title = strings.TrimSpace(a.Title)
	a.Subtitle = strings.TrimSpace(a.Subtitle)
	a.Content = strings.TrimSpace(a.Content)
	var lastInsertID int
	err = db.QueryRow("INSERT INTO info1(title,subtitle,content,creationtimestamp) VALUES($1,$2,$3,$4) returning id;", a.Title, a.Subtitle, a.Content, time.Now()).Scan(&lastInsertID)
	a.CreationTimestamp = time.Now()
	a.ID = int64(lastInsertID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"` + err.Error() + `"`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Article with title %v got posted successfully!", a.Title)))
}

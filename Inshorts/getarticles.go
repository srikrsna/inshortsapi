package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

//GETHandler to get all the articles
func GETHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	//pagination list using limit and offset query parameters
	limit := query.Get("limit")
	offset := query.Get("offset")
	db := OpenConnection()
	defer db.Close()
	var rows *sql.Rows
	var err error
	mutex.Lock()
	defer mutex.Unlock()
	switch {
	case limit == "" && offset != "":
		sqlstatement := "SELECT * FROM info1 ORDER BY creationtimestamp DESC OFFSET $1 "
		rows, err = db.Query(sqlstatement, offset)
	case limit != "" && offset == "":
		sqlstatement := "SELECT * FROM info1 ORDER BY creationtimestamp DESC LIMIT $1 "
		rows, err = db.Query(sqlstatement, limit)
	case limit == "" && offset == "":
		sqlstatement := "SELECT * FROM info1 ORDER BY creationtimestamp DESC"
		rows, err = db.Query(sqlstatement)
	default:
		sqlstatement := "SELECT * FROM info1 ORDER BY creationtimestamp DESC LIMIT $1 OFFSET $2 "
		rows, err = db.Query(sqlstatement, limit, offset)
	}
	defer rows.Close()
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var all []Article
	for rows.Next() {
		var article Article
		rows.Scan(&article.ID, &article.Title, &article.Subtitle, &article.Content, &article.CreationTimestamp)
		all = append(all, article)
	}
	peopleBytes, err := json.MarshalIndent(all, "", "\t")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(peopleBytes)
}

package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
)

//Handler to implement search on title,subtitle,content
func queryHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	q := strings.TrimSpace(query.Get("q"))
	// fmt.Println("TYPE", reflect.TypeOf(q))
	// fmt.Println(q)
	// fmt.Println(template.HTMLEscapeString(q))
	check := template.HTMLEscapeString(q)
	if check != q || len(q) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"Unaccepted/malicious input"} Entered:"` + check + `" "` + q + `"`))
		return
	}
	db := OpenConnection()
	defer db.Close()
	//rows, err := db.Query(`SELECT * FROM info1 WHERE title=$1 OR subtitle=$1 OR content=$1`, q)
	rows, err := db.Query(`SELECT * FROM info1 WHERE title like $1 OR subtitle like $1 OR content like $1`, "%"+q+"%")
	if err != nil {
		log.Fatalf("Error message,%v", err)
	}
	var all []Article
	for rows.Next() {
		var a Article
		err1 := rows.Scan(&a.ID, &a.Title, &a.Subtitle, &a.Content, &a.CreationTimestamp)
		all = append(all, a)
		if err1 != nil {
			log.Fatalf("Unable to fetch article. %v", err1)
		}
	}
	articlequeried, err := json.MarshalIndent(all, "", "\t")
	checkErr(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(articlequeried)
	//fmt.Println("Query exhibited!")
}

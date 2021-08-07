package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

//Article gives model struct
type Article struct {
	ID                int64     `json:"id"`
	Title             string    `json:"title"`
	Subtitle          string    `json:"subtitle"`
	Content           string    `json:"content"`
	CreationTimestamp time.Time `json:"creationtimestamp"`
}

//Postgres connection consts
const (
	user     = "postgres"
	password = "1357"
	dbname   = "postgres"
)

//OpenConnection connects to the db
func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("user=%s "+"password=%s dbname=%s sslmode=disable", user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	checkErr(err)
	err = db.Ping()
	checkErr(err)
	//fmt.Println("Connected!")
	return db
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func multiplexer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GETHandler(w, r)
	case "POST":
		POSTHandler(w, r)
	}
}

func main() {
	db := OpenConnection()
	defer db.Close()
	http.HandleFunc("/articles", multiplexer)
	http.HandleFunc("/articles/", singleGETHandler)
	http.HandleFunc("/articles/search", queryHandler)
	http.HandleFunc("/articles/delete/", deleteHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

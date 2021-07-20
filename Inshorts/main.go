package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

//Article gives model struct
type Article struct {
	ID                int64  `json:"id"`
	Title             string `json:"title"`
	Subtitle          string `json:"subtitle"`
	Content           string `json:"content"`
	CreationTimestamp string `json:"creationtimestamp"`
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
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected!")
	return db
}

//GETHandler to get all the articles
func GETHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	//pagination list using limit and offset query parameters
	limit := query.Get("limit")
	offset := query.Get("offset")
	db := OpenConnection()
	var rows *sql.Rows
	var err error
	switch {
	case limit == "" && offset != "":
		sqlstatement := "SELECT * FROM info OFFSET $1"
		rows, err = db.Query(sqlstatement, offset)
	case limit != "" && offset == "":
		sqlstatement := "SELECT * FROM info LIMIT $1"
		rows, err = db.Query(sqlstatement, limit)
	case limit == "" && offset == "":
		sqlstatement := "SELECT * FROM info"
		rows, err = db.Query(sqlstatement)
	default:
		sqlstatement := "SELECT * FROM info LIMIT $1 OFFSET $2"
		rows, err = db.Query(sqlstatement, limit, offset)
	}
	checkErr(err)
	var all []Article
	for rows.Next() {
		var article Article
		rows.Scan(&article.ID, &article.Title, &article.Subtitle, &article.Content, &article.CreationTimestamp)
		all = append(all, article)
	}
	peopleBytes, _ := json.MarshalIndent(all, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)
	fmt.Println("Got all articles!")
	defer rows.Close()
	defer db.Close()
}

//singleGETHandler to fetch single article based on id
func singleGETHandler(w http.ResponseWriter, r *http.Request) {
	i := r.URL.Path[len("/articles/"):]
	fmt.Println("The ID is " + i)
	id := string(i)
	db := OpenConnection()
	var a Article
	sqlStatement := `SELECT * FROM info WHERE id=$1`
	row := db.QueryRow(sqlStatement, id)
	err1 := row.Scan(&a.ID, &a.Title, &a.Subtitle, &a.Content, &a.CreationTimestamp)
	if err1 != nil {
		log.Fatalf("Unable to get article. %v", err1)
	}
	singlearticle, _ := json.MarshalIndent(a, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(singlearticle)
	fmt.Println("id called:", id)
	defer db.Close()
}

//POSTHandler to insert an article and store in db
func POSTHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	var a Article
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var lastInsertId int
	err = db.QueryRow("INSERT INTO info(title,subtitle,content,creationtimestamp) VALUES($1,$2,$3,$4) returning id;", a.Title, a.Subtitle, a.Content, time.Now()).Scan(&lastInsertId)
	checkErr(err)
	fmt.Println("last inserted id =", lastInsertId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//Handler to implement search on title,subtitle,content
func queryHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	q := query.Get("q")
	fmt.Println(q)
	db := OpenConnection()
	var a Article
	rows, err := db.Query(`SELECT * FROM info WHERE title=$1 OR subtitle=$1 OR content=$1`, q)
	if err != nil {
		log.Fatalf("Not getting query,%v", err)
	}
	for rows.Next() {
		err1 := rows.Scan(&a.ID, &a.Title, &a.Subtitle, &a.Content, &a.CreationTimestamp)
		if err1 != nil {
			log.Fatalf("Unable to get article. %v", err1)
		}
	}
	articlequeried, _ := json.MarshalIndent(a, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(articlequeried)
	fmt.Println("Query exhibited!")
	defer db.Close()
}

//Handler to delete an article wrt id
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	i := r.URL.Path[len("/articles/delete/"):]
	fmt.Println("The deleted ID is " + i)
	id := string(i)
	db := OpenConnection()
	stmt, err := db.Prepare("DELETE from info WHERE id=$1")
	checkErr(err)
	res, err := stmt.Exec(id)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(affect, "rows changed")
	defer db.Close()
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

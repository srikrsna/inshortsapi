package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"
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

var mutex sync.Mutex

//Postgres connection consts
const (
	user     = "postgres"
	password = "1357"
	dbname   = "postgres"
)

// TODO: Always return error back to the caller. Panic should not be used in this situation
// Eg: OpenConnection() (*sql.DB, error) {
// 		db, err := sql.Open(...)
//      if err != nil {
//    		return nil, fmt.Errorf("unable to connect to db: %w", err)
//      }
// 		...
// }
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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`method not allowed`))
	}
}

func main() {
	db := OpenConnection()
	defer db.Close()
	// TODO: Where is this DB being used?
	// These are accurate
	// TODO: Handlers should be named better.
	http.HandleFunc("/articles", multiplexer)
	http.HandleFunc("/articles/", singleGETHandler)
	http.HandleFunc("/articles/search", queryHandler)
	http.HandleFunc("/articles/delete/", deleteHandler)
	// TODO: Err needs to be checked here. Fatal is only required if err is not nil/unexpected one.
	log.Fatal(http.ListenAndServe(":8080", nil))
}

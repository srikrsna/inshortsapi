package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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

// OpenDb
func OpenDb(user, password, dbname string) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping db: %w", err)
	}

	return db, nil
}

func run() error {
	var (
		dbu = flag.String("db.user", "postgres", "Database username")
		dbp = flag.String("db.pass", "1357", "Database password")
		dbn = flag.String("db.name", "postgres", "Database name")
	)

	flag.Parse()

	db, err := OpenDb(*dbu, *dbp, *dbn)
	if err != nil {
		return err
	}
	defer db.Close()

	ah := NewArticleHandler(db)

	mux := http.NewServeMux()
	mux.Handle("/articles", http.StripPrefix("/articles", ah))

	return http.ListenAndServe(":8080", ContentCheck(mux))
}

type ArticleHandler struct {
	db *sql.DB
	http.Handler
}

func NewArticleHandler(db *sql.DB) *ArticleHandler {
	ah := &ArticleHandler{}

	mux := http.NewServeMux()

	mux.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ah.GetArticle(w, r)
		case http.MethodPost:
			ah.CreateArticle(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`method not allowed`))
		}
	})
	mux.HandleFunc("/", ah.GetArticle)
	mux.HandleFunc("/search", ah.SearchArticle)
	mux.HandleFunc("/delete", ah.DeleteArticle)

	ah.Handler = ContentCheck(mux)
	return ah
}

func ContentCheck(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut, http.MethodPost, http.MethodPatch:
			ct := r.Header.Get("content-type")
			if ct != "application/json" {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(fmt.Sprintf("Expectedjson type input got '%s'", ct)))
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}

type Error struct {
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&Error{Message: msg})
}

func (ah *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var a Article
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		log.Println("error decoding json in create article: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(a.Title) == 0 || len(a.Subtitle) == 0 || len(a.Content) == 0 {
		WriteError(w, http.StatusBadRequest, "Title/Subtitle/Content Cannot be empty")
		return
	}

	a.Title = strings.TrimSpace(a.Title)
	a.Subtitle = strings.TrimSpace(a.Subtitle)
	a.Content = strings.TrimSpace(a.Content)
	a.CreationTimestamp = time.Now()

	var id int
	if err := ah.db.QueryRow(
		"INSERT INTO articles(title, subtitle, content, creation_timestamp) VALUES($1, $2, $3, $4) returning id",
		a.Title,
		a.Subtitle,
		a.Content,
		a.CreationTimestamp,
	).Scan(&id); err != nil {
		log.Println("unable to insert article: ", err)
		WriteError(w, http.StatusInternalServerError, "Unknow error")
		return
	}

	a.ID = int64(id)

	WriteResponse(w, a)
}

func (ah *ArticleHandler) ListArticles(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	// pagination list using limit and offset query parameters
	limit := query.Get("limit")
	offset := query.Get("offset")

	args := []interface{}{}
	sql := selBase + " ORDER DESC "
	if limit != "" {
		sql += " LIMIT $1"
		args = append(args, limit)
	}

	if offset != "" {
		sql += " OFFSET $" + strconv.Itoa(len(args)+1)
	}

	ah.listArticles(w, sql, args...)
}

func (ah *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	i := r.URL.Path[len("/articles/"):]
	id, err := strconv.ParseInt(i, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "malformed id, should be int64")
		return
	}

	db := ah.db
	var a Article
	const sql = selBase + ` WHERE id = $1`
	if err := ScanArticle(db.QueryRow(sql, id), &a); err != nil {
		log.Println("unable to get article from db: ", err)
		WriteError(w, http.StatusInternalServerError, "Unknown Error")
		return
	}

	WriteResponse(w, a)
}

func (ah *ArticleHandler) listArticles(w http.ResponseWriter, q string, args ...interface{}) {
	rows, err := ah.db.Query(q, args...)
	if err != nil {
		log.Fatalf("Error message,%v", err)
	}
	defer rows.Close()

	if err != nil {
		log.Println("unable to query articles: ", err)
		WriteError(w, http.StatusInternalServerError, "Unknown")
		return
	}
	defer rows.Close()

	var all []*Article
	for rows.Next() {
		var article Article
		if err := ScanArticle(rows, &article); err != nil {
			log.Println("unable to scan article: ", err)
			WriteError(w, http.StatusInternalServerError, "Unknown")
			return
		}
		all = append(all, &article)
	}

	WriteResponse(w, all)
}

func (ah *ArticleHandler) SearchArticle(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	q := strings.TrimSpace(query.Get("q"))

	ah.listArticles(w, selBase+` WHERE title like $1 OR subtitle like $1 OR content like $1`, "%"+q+"%")
}

func (ah *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	i := r.URL.Path[len("/delete/"):]
	id, err := strconv.ParseInt(i, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "malformed id, should be int64")
		return
	}

	if _, err := ah.db.Exec("DELETE from articles WHERE id = $1", id); err != nil {
		log.Println("unable to insert article: ", err)
		WriteError(w, http.StatusInternalServerError, "Unknown Error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func WriteResponse(w http.ResponseWriter, v interface{}) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("unable to write create article response: %v", err)
	}
}

type Scanner interface {
	Scan(dest ...interface{}) error
}

const selBase = "SELECT id, title, subtitle, content, creation_timestamp from articles"

func ScanArticle(row Scanner, a *Article) error {
	return row.Scan(a.ID, a.Title, a.Subtitle, a.Content, a.CreationTimestamp)
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

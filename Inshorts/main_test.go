package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func getAh(t *testing.T) *ArticleHandler {
	db, err := OpenDb("", "", "")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return NewArticleHandler(db)
}

func TestArticleHandler(t *testing.T) {
	ah := getAh(t)
	if _, err := ah.db.Exec("TRUNCATE articles"); err != nil {
		t.Fatal(err)
	}

	t.Run("Create", func(t *testing.T) {
		exp := Article{
			Title:    "Title 14",
			Subtitle: "Sub 14",
			Content:  "Trending 14",
		}

		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(exp)

		req, err := http.NewRequest("POST", "", &buf)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		r := httptest.NewRecorder()
		ah.ServeHTTP(r, req)
		if status := r.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var act Article
		if err := json.NewDecoder(r.Body).Decode(&act); err != nil {
			t.Fatal("should return article: ", err)
		}

		tdiff := time.Since(act.CreationTimestamp)
		if tdiff > time.Second {
			t.Error("invalid creation timestamp diff", tdiff)
		}

		exp.ID = act.ID
		exp.CreationTimestamp = act.CreationTimestamp
		if reflect.DeepEqual(act, exp) {
			t.Errorf("mismatch: got %v want %v", act, exp)
		}

		t.Run("Get", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/articles/"+strconv.Itoa(int(exp.ID)), nil)
			if err != nil {
				t.Fatal(err)
			}
			r := httptest.NewRecorder()

			ah.ServeHTTP(r, req)
			if status := r.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			var act Article
			if err := json.NewDecoder(r.Body).Decode(&act); err != nil {
				t.Fatal("should return article: ", err)
			}

			if reflect.DeepEqual(act, exp) {
				t.Errorf("mismatch: got %v want %v", act, exp)
			}
		})

		// Delete and Get
	})
}

func TestGETHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/articles?offset=1&limit=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if req == nil {
		t.Fatal("Null returned")
	}
}

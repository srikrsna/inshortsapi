package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestsingleGETHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/articles/4", nil)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRecorder()
	handler := http.HandlerFunc(singleGETHandler)
	handler.ServeHTTP(r, req)
	if status := r.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := `{"id":4,"title":"Title 2","subtitle":"Sub 2","content":"Trending 2","creationtimestamp":"0000-01-01T19:23:38.171134Z"}`
	if r.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			r.Body.String(), expected)
	}
}
func TestGETHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/articles", nil)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRecorder()
	handler := http.HandlerFunc(GETHandler)
	handler.ServeHTTP(r, req)
	if status := r.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := `[{"id":3,"title":"Title 1","subtitle":"Sub 1","content":"Trending 1","creationtimestamp":"0000-01-01T19:23:29.764641Z"},{"id":4,"title":"Title 2","subtitle":"Sub 2","content":"Trending 2","creationtimestamp":"0000-01-01T19:23:38.171134Z"},{"id":5,"title":"Title 3","subtitle":"Sub 3","content":"Trending 3","creationtimestamp":"0000-01-01T19:23:47.57623Z"},{"id":6,"title":"Title 4","subtitle":"Sub 4","content":"Trending 4","creationtimestamp":"0000-01-01T19:23:57.31599Z"},{"id":7,"title":"Title 5","subtitle":"Sub 5","content":"Trending 5","creationtimestamp":"0000-01-01T19:24:05.210796Z"},{"id":8,"title":"Title 6","subtitle":"Sub 6","content":"Trending 6","creationtimestamp":"0000-01-01T19:24:14.144641Z"},{"id":9,"title":"Title 7","subtitle":"Sub 7","content":"Trending 7","creationtimestamp":"0000-01-01T19:24:23.7952Z"},{"id":10,"title":"Title 8","subtitle":"Sub 8","content":"Trending 8","creationtimestamp":"0000-01-01T19:24:33.569361Z"},{"id":11,"title":"Title 9","subtitle":"Sub 9","content":"Trending 9","creationtimestamp":"0000-01-01T19:24:42.433275Z"},{"id":13,"title":"Title 11","subtitle":"Sub 11","content":"Trending 11","creationtimestamp":"0000-01-01T19:50:12.68611Z"},{"id":14,"title":"Title 12","subtitle":"Sub 12","content":"Trending 12","creationtimestamp":"0000-01-01T19:42:16.250215Z"}]`
	if r.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			r.Body.String(), expected)
	}
}

func TestPOSTHandler(t *testing.T) {
	var input = []byte(`{"title":"Title 14","subtitle":"Sub 14","content":"Trending 14"}`)

	req, err := http.NewRequest("POST", "/entry", bytes.NewBuffer(input))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(POSTHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := `{0 Title 14 Sub 14 Trending 14 }`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
func TestqueryHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/articles/search", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("q", "Title 2")
	req.URL.RawQuery = q.Encode()
	r := httptest.NewRecorder()
	handler := http.HandlerFunc(queryHandler)
	handler.ServeHTTP(r, req)
	if status := r.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"id":4,"title":"Title 2","subtitle":"Sub 2","content":"Trending 2","creationtimestamp":"0000-01-01T19:23:38.171134Z"}`
	if r.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			r.Body.String(), expected)
	}

}

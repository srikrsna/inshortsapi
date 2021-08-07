package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSingleGETHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/articles/5", nil)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRecorder()
	handler := http.HandlerFunc(singleGETHandler)
	handler.ServeHTTP(r, req)
	if status := r.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := `{"id":5,"title":"Title 3","subtitle":"Sub 3","content":"Trending 3","creationtimestamp":"2021-08-06T23:58:02.143318+05:30"}`
	if r.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			r.Body.String(), expected)
	}
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

func TestPOSTHandler(t *testing.T) {
	var input = []byte(`{"title":"Title 14","subtitle":"Sub 14","content":"Trending 14"}`)
	req, err := http.NewRequest("POST", "/entry", bytes.NewBuffer(input))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	r := httptest.NewRecorder()
	handler := http.HandlerFunc(POSTHandler)
	handler.ServeHTTP(r, req)
	if status := r.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := `Article with title Title 14 got posted successfully!`
	if r.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			r.Body.String(), expected)
	}
}

func TestQueryHandler(t *testing.T) {
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
}
func BenchmarkGETHandler(b *testing.B) {
	for n := 0; n < b.N; n++ {
		http.Get("http://localhost:12345/articles")
	}
}

func BenchmarkSingleGETHandler(b *testing.B) {
	for n := 0; n < b.N; n++ {
		http.Get("http://localhost:12345/articles/2")
	}
}

func BenchmarkQueryHandler(b *testing.B) {
	for n := 0; n < b.N; n++ {
		http.Get("http://localhost:8080/articles/search?q=u")
	}
}

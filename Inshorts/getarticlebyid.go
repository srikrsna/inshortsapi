package main

import (
	"encoding/json"
	"net/http"
)

//singleGETHandler to fetch single article based on id
func singleGETHandler(w http.ResponseWriter, r *http.Request) {
	i := r.URL.Path[len("/articles/"):]
	id := string(i)
	db := OpenConnection()
	defer db.Close()
	mutex.Lock()
	defer mutex.Unlock()
	var a Article
	if len(id) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	sqlStatement := `SELECT * FROM info1 WHERE id=$1`
	row := db.QueryRow(sqlStatement, id)
	err1 := row.Scan(&a.ID, &a.Title, &a.Subtitle, &a.Content, &a.CreationTimestamp)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "message": "` + err1.Error() + `" }`))
		return
	}
	singlearticle, err := json.Marshal(a)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(singlearticle)

}

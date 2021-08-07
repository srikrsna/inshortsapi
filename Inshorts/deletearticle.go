package main

import (
	"fmt"
	"net/http"
)

//Handler to delete an article wrt id
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	i := r.URL.Path[len("/articles/delete/"):]
	//fmt.Println("The deleted ID is " + i)
	id := string(i)
	db := OpenConnection()
	defer db.Close()
	stmt, err := db.Prepare("DELETE from info1 WHERE id=$1")
	checkErr(err)
	res, err := stmt.Exec(id)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(affect, "rows changed")
}

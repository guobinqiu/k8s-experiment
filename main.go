package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func main() {
	fmt.Println("Go Web App Started on Port 3001")
	setupRoutes()
	http.ListenAndServe(":3001", nil)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "My Awesome Go App!!!")
}

func userPage(w http.ResponseWriter, r *http.Request) {
	db := getDB()
	rows, _ := db.Query("SELECT * FROM user")
	for rows.Next() {
		var name string
		var age int
		rows.Scan(&name, &age)
		fmt.Fprintf(w, "name=%s, age=%d\n", name, age)
	}
}

func getDB() *sql.DB {
	db, err := sql.Open("mysql", "guobin:222222@tcp(mysql:3306)/test?charset=utf8")
	if err != nil {
		panic(err)
	}
	return db
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/users", userPage)
}

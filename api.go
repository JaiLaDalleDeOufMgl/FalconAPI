package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net"
	"net/http"
)

type Tag struct {
	EXPIRES string `json:"expires"`
}

func Check(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"version" : "0.0.1"}`)
}

func License(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:gghhtt@tcp(127.0.0.1:3306)/falconapi")
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	results, err := db.Query("SELECT expires FROM license WHERE license_key = ? AND ip = ?", r.URL.Path[9:], ip)
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		var tag Tag
		err = results.Scan(&tag.EXPIRES)
		if err != nil {
			panic(err.Error())
		}
		fmt.Fprintf(w, `{status: "success", expires: "%s"}`, tag.EXPIRES)
	}
}

func main() {
	http.HandleFunc("/check", Check)
	http.HandleFunc("/license/", License)
	http.ListenAndServe(":8080", nil)
}

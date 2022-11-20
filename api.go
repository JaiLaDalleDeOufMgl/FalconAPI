package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net"
	"net/http"
	"time"
)

type Tag struct {
	DAY   string `json:"day"`
	MONTH string `json:"month"`
	YEAR  string `json:"year"`
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
	results, err := db.Query("SELECT day, month, year FROM license WHERE license_key = ? AND ip = ?", r.URL.Path[9:], ip)
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		var tag Tag
		err = results.Scan(&tag.DAY, &tag.MONTH, &tag.YEAR)
		if err != nil {
			panic(err.Error())
		}
		dt := time.Now()
		if tag.YEAR >= dt.Format("2006") {
			if tag.MONTH >= dt.Format("01") {
				if tag.DAY >= dt.Format("02") {
					fmt.Println("license ok")
					fmt.Fprintf(w, `{status: "success"}`)
				}
			}
		}
	}
}

func main() {
	http.HandleFunc("/check", Check)
	http.HandleFunc("/license/", License)
	http.ListenAndServe(":8080", nil)
}

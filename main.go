package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"

	_ "github.com/lib/pq"
)

var m sync.Mutex

var data = make(map[int]string)

const MIN_RANGE = 100000
const MAX_RANGE = 999999

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	UID := rand.Intn(MAX_RANGE-MIN_RANGE) + MIN_RANGE
	m.Lock()
	defer m.Unlock()
	_, exists := data[UID]
	tmp := UID
	for exists {
		UID++
		if UID > MAX_RANGE {
			break
		}
		_, exists = data[UID]
	}
	UID = tmp
	for exists {
		UID--
		if UID < MIN_RANGE {
			fmt.Fprint(w, "No empty slots available.")
			return
		}
		_, exists = data[UID]
	}
	insertDb(UID, url)
	data[UID] = url
	fmt.Fprint(w, UID)
}

func shortnerRouter(w http.ResponseWriter, r *http.Request) {
	m.Lock()
	defer m.Unlock()
	id, err := strconv.Atoi((r.URL.Query().Get("id")))
	if err != nil {
		fmt.Fprint(w, "Bad Path")
		return
	}
	val, exists := data[id]
	if !exists {
		fmt.Fprint(w, "Bad Path")
		return
	}
	http.Redirect(w, r, val, http.StatusSeeOther)
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "12345678"
	dbname   = "urldb"
)

func getConnStr() string {
	return fmt.Sprintf("host = %s port = %d user = %s password = %s dbname = %s sslmode=disable", host, port, user, password, dbname)
}

func insertDb(UID int, url string) {
	db, err := sql.Open("postgres", getConnStr())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("Insert into links(uid, url) values($1, $2)", UID, url)
	if err != nil {
		log.Fatal(err)
	}
}

func initDb() {
	db, err := sql.Open("postgres", getConnStr())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT uid, url from links")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var UID int
		var url string
		err = rows.Scan(&UID, &url)
		if err != nil {
			panic(err)
		}
		data[UID] = url
	}
}

func main() {
	initDb()
	router := http.NewServeMux()
	router.HandleFunc("POST /shorten/", shortenHandler)
	router.HandleFunc("GET /get/", shortnerRouter)
	http.ListenAndServe(":8080", router)
}

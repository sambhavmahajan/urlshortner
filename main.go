package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

var m sync.Mutex

var data = make(map[int]string)

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	UID := rand.Intn(999999-100000) + 100000
	m.Lock()
	defer m.Unlock()
	_, exists := data[UID]
	tmp := UID
	for exists {
		UID++
		if UID > 999999 {
			break
		}
		_, exists = data[UID]
	}
	UID = tmp
	for exists {
		UID++
		if UID < 100000 {
			fmt.Fprint(w, "No empty slots available.")
			return
		}
		_, exists = data[UID]
	}
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

func main() {
	router := http.NewServeMux()
	router.HandleFunc("POST /shorten/", shortenHandler)
	router.HandleFunc("GET /get/", shortnerRouter)
	http.ListenAndServe(":8080", router)
}

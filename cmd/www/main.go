package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	verbose = false
	addr    = flag.String("addr", "localhost:8000", "serve address")
)

func handleApi(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/")
	if verbose {
		fmt.Println(r.Method, name)
	}

	if r.Method == "GET" {

	} else {
		w.WriteHeader(404)
	}
}

func main() {
	flag.Parse()

	server := http.FileServer(http.Dir("www"))
	http.Handle("/", server)
	http.HandleFunc("/api/", handleApi)

	for {
		err := http.ListenAndServe(*addr, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

package main

import (
	"flag"
	"log"

	"trib/store"
	"triblab"
)

var (
	addr = flag.String("addr", "localhost:9000", "backend serve address")
)

func main() {
	flag.Parse()

	s := store.NewStorage()

	e := triblab.ServeBack(*addr, s, nil)
	if e != nil {
		log.Fatal(e)
	}
}

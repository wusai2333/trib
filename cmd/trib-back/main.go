// Tribbler backend launcher program.
package main

import (
	"flag"
	"log"

	"trib/randaddr"
	"trib/store"
	"triblab"
)

var (
	addr = flag.String("addr", "localhost:rand", "backend serve address")
)

func main() {
	flag.Parse()

	*addr = randaddr.Resolve(*addr)

	s := store.NewStorage()

	log.Printf("tribble backend serve on %s", *addr)

	e := triblab.ServeBack(*addr, s, nil)
	if e != nil {
		log.Fatal(e)
	}
}

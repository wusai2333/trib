package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"trib/store"
	"triblab"
)

var (
	addr = flag.String("addr", "localhost:rand", "backend serve address")
)

func main() {
	flag.Parse()

	if strings.HasSuffix(*addr, ":rand") {
		*addr = strings.TrimSuffix(*addr, ":rand")
		rand.Seed(time.Now().UnixNano())
		port := 10000 + int(rand.Uint32()%20000)
		*addr = fmt.Sprintf("%s:%d", *addr, port)
	}

	s := store.NewStorage()

	log.Printf("tribble backend serve on %s", *addr)

	e := triblab.ServeBack(*addr, s, nil)
	if e != nil {
		log.Fatal(e)
	}
}

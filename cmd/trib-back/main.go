// Tribbler backend launcher program.
package main

import (
	"flag"
	"fmt"
	"log"

	"trib"
	"trib/store"
	"triblab"
)

var (
	frc   = flag.String("rc", "trib.rc", "tribbler service config")
	index = flag.Int("index", 0, "index in the back-end list")
)

func noError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	flag.Parse()

	rc, e := trib.LoadRC(*frc)
	noError(e)

	if *index >= rc.BackCount() {
		e = fmt.Errorf("back-end index out of range: %d/%d",
			*index, rc.BackCount())
		log.Fatal(e)
	}

	backConfig := rc.BackConfig(*index, store.NewStorage())

	log.Printf("tribble backend serve on %s, peer on %s",
		backConfig.Addr, backConfig.Peer.Addr(),
	)

	noError(triblab.ServeBack(backConfig))
}

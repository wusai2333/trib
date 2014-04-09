// Tribbler back-end launcher.
package main

import (
	"flag"
	"log"
	"strconv"

	"trib"
	"trib/local"
	"trib/store"
	"triblab"
)

var (
	frc = flag.String("rc", "trib.rc", "tribbler service config file")
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

	run := func(i int) {
		backConfig := rc.BackConfig(i, store.NewStorage())
		log.Printf("tribbler back-end serve on %s", backConfig.Addr)
		noError(triblab.ServeBack(backConfig))
	}

	args := flag.Args()

	if len(args) == 0 {
		// scan for addresses on this machine
		n := 0
		for i, b := range rc.Backs {
			if local.Check(b) {
				go run(i)
				n++
			}
		}

		if n == 0 {
			log.Fatal("no back-end found for this host")
		}
	} else {
		// scan for indices for the addresses
		for _, a := range args {
			i, e := strconv.Atoi(a)
			noError(e)
			go run(i)
		}
	}

	select {}
}

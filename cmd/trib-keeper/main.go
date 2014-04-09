// Tribbler back-end keeper launcher.
package main

import (
	"flag"
	"log"
	"strconv"

	"trib"
	"trib/local"
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
		keeperConfig := rc.KeeperConfig(i)
		log.Printf("tribbler keeper serve on %s", keeperConfig.Addr())
		noError(triblab.ServeKeeper(keeperConfig))
	}

	args := flag.Args()
	if len(args) == 0 {
		n := 0
		for i, k := range rc.Keepers {
			if local.Check(k) {
				go run(i)
				n++
			}
		}

		if n == 0 {
			log.Fatal("no keeper found for this host")
		}
	} else {
		for _, a := range args {
			i, e := strconv.Atoi(a)
			noError(e)

			go run(i)
		}
	}

	select {}
}

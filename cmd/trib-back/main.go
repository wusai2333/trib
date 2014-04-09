// Tribbler backend launcher program.
package main

import (
	"flag"
	"log"
	"net"
	"strconv"
	"strings"

	"trib"
	"trib/store"
	"triblab"
)

var (
	frc = flag.String("rc", "trib.rc", "tribbler service config")
)

func noError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

var localAddrs = func() []string {
	addrs, e := net.InterfaceAddrs()
	noError(e)

	ret := make([]string, 0, len(addrs))

	for _, addr := range addrs {
		ret = append(ret, addr.String())
	}

	return ret
}()

func onThisMachine(addr string) bool {
	a, e := net.ResolveTCPAddr("tcp", addr)
	noError(e)

	ip := a.IP.String()
	for _, addr := range localAddrs {
		if strings.HasPrefix(addr, ip) {
			return true
		}
	}

	return false
}

func main() {
	flag.Parse()

	rc, e := trib.LoadRC(*frc)
	noError(e)

	run := func(i int) {
		backConfig := rc.BackConfig(i, store.NewStorage())
		log.Printf("tribble back-end serve on %s", backConfig.Addr)
		noError(triblab.ServeBack(backConfig))
	}

	args := flag.Args()

	if len(args) == 0 {
		// scan for addresses on this machine
		n := 0
		for i, b := range rc.Backs {
			if onThisMachine(b) {
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
			if e != nil {
				log.Fatal(e)
			}
			go run(i)
		}
	}

	select {}
}

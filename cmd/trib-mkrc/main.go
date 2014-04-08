package main

import (
	"flag"
	"fmt"
	"log"

	"trib"
	"trib/randaddr"
)

var (
	local = flag.Bool("local", false, "always use local ports")
	n     = flag.Int("n", 1, "number of servers")
)

func main() {
	if *n > 100 {
		log.Fatal(fmt.Errorf("too many servers"))
	}

	p := randaddr.RandPort()

	rc := new(trib.RC)
	rc.Backs = make([]*trib.BackAddr, *n)

	if !*local {
		const ipOffset = 211
		for i := 0; i < *n; i++ {
			host := fmt.Sprintf("172.22.14.%d", ipOffset+i)
			saddr := fmt.Sprintf("%s:%d", host, p)
			paddr := fmt.Sprintf("%s:%d", host, p+1)
			rc.Backs[i] = &trib.BackAddr{saddr, paddr}

			p += 2
		}
	} else {
		// TODO:
	}
}

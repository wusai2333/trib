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
	frc   = flag.String("rc", "trib.rc", "back-end config file")
)

func main() {
	flag.Parse()

	if *n > 10 {
		log.Fatal(fmt.Errorf("too many servers"))
	}

	p := randaddr.RandPort()

	rc := new(trib.RC)
	rc.Backs = make([]string, *n)

	if !*local {
		const ipOffset = 211
		for i := 0; i < *n; i++ {
			host := fmt.Sprintf("172.22.14.%d", ipOffset+i)
			rc.Backs[i] = fmt.Sprintf("%s:%d", host, p)
		}
	} else {
		for i := 0; i < *n; i++ {
			rc.Backs[i] = fmt.Sprintf("localhost:%d", p+i)
		}
	}

	rc.Keepers = make([]string, 0, 3)

	fmt.Println(rc.String())

	if *frc != "" {
		e := rc.Save(*frc)
		if e != nil {
			log.Fatal(e)
		}
	}
}

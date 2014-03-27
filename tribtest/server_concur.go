package tribtest

import (
	"runtime/debug"
	"strconv"
	"testing"
	"time"

	"trib"
)

func CheckServerConcur(t *testing.T, server trib.Server) {
	ne := func(e error) {
		if e != nil {
			debug.PrintStack()
			t.Fatal(e)
		}
	}

	as := func(cond bool) {
		if !cond {
			debug.PrintStack()
			t.Fatal()
		}
	}

	ne(server.SignUp("user"))
	tm := time.Now()

	p := func(th, n int, done chan<- bool) {
		for i := 0; i < n; i++ {
			ne(server.Post("user", "", strconv.Itoa(th*100+n), tm))
		}
		done <- true
	}

	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go p(i, 10, done)
	}

	for i := 0; i < 5; i++ {
		<-done
	}

	ret, e := server.Tribs("user")
	ne(e)
	as(len(ret) == 50)
}

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"trib"
	"triblab"
)

func noError(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}

func kv(k, v string) *trib.KeyValue {
	return &trib.KeyValue{k, v}
}

func kva(args []string) *trib.KeyValue {
	return kv(args[2], args[3])
}

func pata(args []string) *trib.Pattern {
	if len(args) == 2 {
		return pat("", "")
	} else if len(args) == 3 {
		return pat(args[2], "")
	}

	return pat(args[2], args[3])
}

func pat(pre, suf string) *trib.Pattern {
	return &trib.Pattern{pre, suf}
}

func printList(lst trib.List) {
	for _, e := range lst.L {
		fmt.Println(e)
	}
}

func main() {
	flag.Parse()
	args := flag.Args()

	s := triblab.NewClient(args[0])

	var v string
	var b bool
	var lst trib.List
	var n int
	var cret uint64

	switch args[1] {
	case "get":
		noError(s.Get(args[2], &v))
		fmt.Println(v)
	case "set":
		noError(s.Set(kva(args), &b))
		fmt.Println(b)
	case "keys":
		noError(s.Keys(pata(args), &lst))
		printList(lst)
	case "list-get":
		noError(s.ListGet(args[2], &lst))
		printList(lst)
	case "list-append":
		noError(s.ListAppend(kva(args), &b))
		fmt.Println(b)
	case "list-remove":
		noError(s.ListRemove(kva(args), &n))
		fmt.Println(n)
	case "list-keys":
		noError(s.ListKeys(pata(args), &lst))
		printList(lst)
	case "clock":
		var c uint64
		var e error
		if len(args) >= 3 {
			c, e = strconv.ParseUint(args[2], 10, 64)
			noError(e)
		}
		noError(s.Clock(c, &cret))
		fmt.Println(cret)
	}

}

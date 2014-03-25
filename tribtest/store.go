package tribtest

import (
	"sort"
	"testing"

	"trib"
)

func CheckStorage(t *testing.T, s trib.Storage) {
	var v string
	var b bool
	var l = new(trib.List)

	ne := func(e error) {
		if e != nil {
			t.Fatal(e)
		}
	}

	as := func(cond bool) {
		if !cond {
			t.Fail()
		}
	}

	kv := func(k, v string) *trib.KeyValue {
		return &trib.KeyValue{k, v}
	}

	pat := func(pre, suf string) *trib.Pattern {
		return &trib.Pattern{pre, suf}
	}

	v = "_"
	ne(s.Get("", &v))
	as(v == "")

	v = "_"
	ne(s.Get("hello", &v))
	as(v == "")

	ne(s.Set(kv("h8liu", "run"), &b))
	as(b)
	v = ""
	ne(s.Get("h8liu", &v))
	as(v == "run")

	ne(s.Set(kv("h8liu", "Run"), &b))
	as(b)
	v = ""
	ne(s.Get("h8liu", &v))
	as(v == "Run")

	ne(s.Set(kv("h8liu", ""), &b))
	as(b)
	v = "_"
	ne(s.Get("h8liu", &v))
	as(v == "")

	ne(s.Set(kv("h8liu", "k"), &b))
	as(b)
	v = "_"
	ne(s.Get("h8liu", &v))
	as(v == "k")

	ne(s.Set(kv("h8he", "something"), &b))
	as(b)
	v = "_"
	ne(s.Get("h8he", &v))
	as(v == "something")

	ne(s.Keys(pat("h8", ""), l))
	sort.Strings(l.L)
	as(len(l.L) == 2)
	as(l.L[0] == "h8he")
	as(l.L[1] == "h8liu")
}

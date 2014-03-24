package store_test

import (
	"testing"

	"trib"
	. "trib/store"
)

func TestStorage(t *testing.T) {
	s := NewStorage()

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

	var v string
	var b bool

	kv := func(k, v string) *trib.KeyValue {
		return &trib.KeyValue{k, v}
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

	ne(s.Append(kv("h8liu", "ner"), &b))
	as(b)
	v = ""
	ne(s.Get("h8liu", &v))
	as(v == "Runner")

	ne(s.Set(kv("h8liu", ""), &b))
	as(b)
	v = "_"
	ne(s.Get("h8liu", &v))
	as(v == "")

	ne(s.Append(kv("h8liu", "ner"), &b))
	as(b)
	v = ""
	ne(s.Get("h8liu", &v))
	as(v == "ner")
}

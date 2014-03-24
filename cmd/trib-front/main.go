package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"trib"
	"trib/ref"
)

var (
	verbose = false
	addr    = flag.String("addr", "localhost:8000", "serve address")
	server  trib.Server
)

func handleApi(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/")
	if verbose {
		fmt.Println(r.Method, name)
	}

	reply := func(obj interface{}) {
		bytes, e := json.Marshal(obj)
		noError(e)

		_, e = w.Write(bytes)
		logError(e)
	}

	bytes, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Println(e)
		return
	}
	input := string(bytes)

	log.Println(name, input)

	switch name {
	case "add-user":
		e = server.SignUp(input)
		if e != nil {
			reply(NewUserList(nil, e))
			break
		}
		ret, e := server.ListUsers()
		reply(NewUserList(ret, e))

	case "list-users":
		ret, e := server.ListUsers()
		reply(NewUserList(ret, e))

	case "list-tribs":
		tribs, e := server.Tribs(input)
		reply(NewTribList(tribs, e))

	case "list-home":
		tribs, e := server.Home(input)
		reply(NewTribList(tribs, e))

	case "is-following":
		ww := new(WhoWhom)
		e := json.Unmarshal(bytes, ww)
		if e != nil {
			reply(NewBool(false, e))
			break
		}
		ret, e := server.IsFollowing(ww.Who, ww.Whom)
		reply(NewBool(ret, e))

	case "follow":
		ww := new(WhoWhom)
		e := json.Unmarshal(bytes, ww)
		if e != nil {
			reply(NewBool(false, e))
			break
		}
		e = server.Follow(ww.Who, ww.Whom)
		reply(NewBool(e == nil, e))

	case "unfollow":
		ww := new(WhoWhom)
		e := json.Unmarshal(bytes, ww)
		if e != nil {
			reply(NewBool(false, e))
			break
		}
		e = server.Unfollow(ww.Who, ww.Whom)
		reply(NewBool(false, e))

	case "post":
		p := new(Post)
		e := json.Unmarshal(bytes, p)
		if e != nil {
			reply(NewBool(false, e))
			break
		}
		e = server.Post(p.Who, p.At, p.Message, time.Now())
		reply(NewBool(e == nil, e))

	default:
		w.WriteHeader(404)
	}
}

func makeServer() trib.Server {
	return ref.NewServer()
}

func populate(server trib.Server) {
	ne := func(e error) {
		if e != nil {
			log.Fatal(e)
		}
	}

	ne(server.SignUp("h8liu"))
	ne(server.SignUp("fenglu"))
	ne(server.SignUp("rkapoor"))

	ne(server.Post("h8liu", "", "Hello, world.", time.Now()))
	ne(server.Post("h8liu", "", "Just tribble it.", time.Now()))
	ne(server.Post("fenglu", "h8liu", "Double tribble.", time.Now()))
	ne(server.Post("rkapoor", "fenglu", "Triple tribble.", time.Now()))

	ne(server.Follow("fenglu", "h8liu"))
	ne(server.Follow("fenglu", "rkapoor"))

	ne(server.Follow("rkapoor", "h8liu"))
}

func main() {
	flag.Parse()
	server = makeServer()
	populate(server)

	http.Handle("/", http.FileServer(http.Dir("www")))
	http.HandleFunc("/api/", handleApi)

	for {
		err := http.ListenAndServe(*addr, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

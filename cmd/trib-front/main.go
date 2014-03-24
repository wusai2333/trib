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

type UserList struct {
	Err   string
	Users []string
}

type TribList struct {
	Err   string
	Tribs []*trib.Trib
}

type Bool struct {
	Err string
	V   bool
}

type WhoWhom struct {
	Who  string
	Whom string
}

func errString(e error) string {
	if e == nil {
		return ""
	}
	return e.Error()
}

func NewTribList(tribs []*trib.Trib, e error) *TribList {
	return &TribList{errString(e), tribs}
}

func NewUserList(users []string, e error) *UserList {
	return &UserList{errString(e), users}
}

func NewBool(b bool, e error) *Bool {
	return &Bool{errString(e), b}
}

func noError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func logError(e error) {
	if e != nil {
		log.Print(e)
	}
}

func handleApi(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/")
	if verbose {
		fmt.Println(r.Method, name)
	}

	jsonReply := func(obj interface{}) {
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
			jsonReply(NewUserList(nil, e))
			break
		}
		ret, e := server.ListUsers()
		jsonReply(NewUserList(ret, e))

	case "list-users":
		ret, e := server.ListUsers()
		jsonReply(NewUserList(ret, e))

	case "list-tribs":
		tribs, e := server.Tribs(input)
		jsonReply(NewTribList(tribs, e))

	case "list-home":
		tribs, e := server.Home(input)
		jsonReply(NewTribList(tribs, e))

	case "is-following":
		ww := new(WhoWhom)
		e := json.Unmarshal(bytes, ww)
		if e != nil {
			jsonReply(NewBool(false, e))
			break
		}
		ret, e := server.IsFollowing(ww.Who, ww.Whom)
		jsonReply(NewBool(ret, e))

	case "follow":
		ww := new(WhoWhom)
		e := json.Unmarshal(bytes, ww)
		if e != nil {
			jsonReply(NewBool(false, e))
			break
		}
		e = server.Follow(ww.Who, ww.Whom)
		if e != nil {
			jsonReply(NewBool(false, e))
			break
		}
		jsonReply(NewBool(true, nil))

	case "unfollow":
		ww := new(WhoWhom)
		e := json.Unmarshal(bytes, ww)
		if e != nil {
			jsonReply(NewBool(false, e))
			break
		}
		e = server.Unfollow(ww.Who, ww.Whom)
		if e != nil {
			jsonReply(NewBool(false, e))
			break
		}
		jsonReply(NewBool(false, nil))

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

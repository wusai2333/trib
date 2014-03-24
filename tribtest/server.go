package tribtest

import (
	"runtime/debug"
	"time"
	"trib"

	"testing"
)

func CheckServer(server trib.Server, t *testing.T) {
	ne := func(e error) {
		if e != nil {
			debug.PrintStack()
			t.Fatal(e)
		}
	}

	er := func(e error) {
		if e == nil {
			debug.PrintStack()
			t.Fatal(e)
		}
	}

	as := func(cond bool) {
		if !cond {
			debug.PrintStack()
			t.Fail()
		}
	}

	ne(server.SignUp("h8liu"))
	er(server.SignUp(" h8liu"))
	er(server.SignUp("8hliu"))
	er(server.SignUp("H8liu"))

	ne(server.SignUp("fenglu"))

	users, e := server.ListUsers()
	ne(e)

	as(len(users) == 2)
	as(users[1] == "h8liu")
	as(users[0] == "fenglu")

	ne(server.Follow("h8liu", "fenglu"))
	b, e := server.IsFollowing("h8liu", "fenglu")
	ne(e)
	as(b)

	b, e = server.IsFollowing("fenglu", "h8liu")
	ne(e)
	as(!b)

	b, e = server.IsFollowing("h8liu", "fenglu2")
	er(e)
	as(!b)

	ne(server.Unfollow("h8liu", "fenglu"))
	er(server.Unfollow("h8liu", "fenglu"))

	b, e = server.IsFollowing("h8liu", "fenglu")
	ne(e)
	as(!b)

	ne(server.Follow("h8liu", "fenglu"))

	tm := time.Now()

	er(server.Post("", "", "", tm))
	er(server.Post("h8liu", "h9liu", "something", tm))

	longMsg := ""
	for i := 0; i < 200; i++ {
		longMsg += " "
	}

	er(server.Post("h8liu", "", longMsg, tm))
	ne(server.Post("h8liu", "", "hello, world", tm))

	tribs, e := server.Tribs("h8liu")
	ne(e)
	as(len(tribs) == 1)
	tr := tribs[0]
	as(tr.User == "h8liu")
	as(tr.Message == "hello, world")
	as(tr.Time == tm)

	tribs, e = server.Home("fenglu")
	ne(e)
	as(len(tribs) == 0)

	ne(server.Follow("fenglu", "h8liu"))
	tribs, e = server.Home("fenglu")
	ne(e)
	as(len(tribs) == 1)
	tr = tribs[0]
	as(tr.User == "h8liu")
	as(tr.Message == "hello, world")
	as(tr.Time == tm)

	tm2 := tm.Add(time.Second)

	ne(server.Post("h8liu", "", "hello, world2", tm2))
	tribs, e = server.Home("fenglu")
	ne(e)
	as(len(tribs) == 2)
	tr = tribs[0]
	as(tr.User == "h8liu")
	as(tr.Message == "hello, world")
	as(tr.Time == tm)

	tr = tribs[1]
	as(tr.User == "h8liu")
	as(tr.Message == "hello, world2")
	as(tr.Time == tm2)
}

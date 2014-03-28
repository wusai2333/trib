package tribtest

import (
	"runtime/debug"
	"sort"
	"testing"
	"time"

	"trib"
)

func CheckServer(t *testing.T, server trib.Server) {
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
			t.Fatal()
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
	sort.Strings(users)
	as(users[0] == "fenglu")
	as(users[1] == "h8liu")

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

	er(server.Post("", "", tm))

	longMsg := ""
	for i := 0; i < 200; i++ {
		longMsg += " "
	}

	er(server.Post("h8liu", longMsg, tm))
	ne(server.Post("h8liu", "hello, world", tm))

	tribs, e := server.Tribs("h8liu")
	ne(e)
	as(len(tribs) == 1)
	tr := tribs[0]
	as(tr.User == "h8liu")
	as(tr.Message == "hello, world")
	as(tr.Time == tm)

	tribs, e = server.Home("fenglu")
	ne(e)
	as(tribs != nil)
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

	ne(server.Post("h8liu", "hello, world2", tm2))
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

	er(server.Follow("fenglu", "fenglu"))
	er(server.Follow("fengl", "fenglu"))
	er(server.Follow("fenglu", "fengl"))
	er(server.Follow("fenglu", "h8liu"))

	tribs, e = server.Home("h8liu")
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

	ne(server.SignUp("rkapoor"))
	fos, e := server.Following("rkapoor")
	ne(e)
	as(fos != nil)
	as(len(fos) == 0)
}

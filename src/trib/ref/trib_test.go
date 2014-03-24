package ref

import (
	"testing"
	"time"
)

func TestTrib(t *testing.T) {
	server := NewServer()
	testServer(server, t)
}

func testServer(server *Server, t *testing.T) {
	ne := func(e error) {
		if e != nil {
			t.Fatal(e)
		}
	}

	er := func(e error) {
		if e == nil {
			t.Fatal(e)
		}
	}

	as := func(cond bool) {
		if !cond {
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
	as(users[0] == "h8liu")
	as(users[1] == "fenglu")

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
}

package ref

import (
	"trib"
)

type following []string

type user struct {
	following
	tribs    []*trib.Trib
	timeline []*trib.Trib
}

type Server struct {
	users map[string]*user
}

var _ trib.Server = new(Server)

func NewServer() *Server {
	ret := &Server{make(map[string]*user)}
	return ret
}

func (self *Server) Register(user string) error {
	panic("todo")
}

func (self *Server) Subscribe(who, whom string) error {
	panic("todo")
}

func (self *Server) Unsubscribe(who, whom string) error {
	panic("todo")
}

func (self *Server) Post(user, post string) error {
	panic("todo")
}

func (self *Server) List(user string, off, count int) ([]*trib.Trib, error) {
	panic("todo")
}

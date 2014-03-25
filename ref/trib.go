package ref

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"trib"
)

type Server struct {
	users map[string]*user
	lock  sync.Mutex
	seq   int
}

var _ trib.Server = new(Server)

func NewServer() *Server {
	ret := &Server{
		users: make(map[string]*user),
	}
	return ret
}

func (self *Server) findUser(user string) (*user, error) {
	ret, found := self.users[user]
	if !found {
		return nil, fmt.Errorf("user %q not exists", user)
	}
	return ret, nil
}

func (self *Server) SignUp(user string) error {
	if len(user) > trib.MaxUsernameLen {
		return fmt.Errorf("username %q too long", user)
	}

	if !trib.IsValidUsername(user) {
		return fmt.Errorf("invalid username %q", user)
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	_, found := self.users[user]
	if found {
		return fmt.Errorf("user %q already exists", user)
	}

	self.users[user] = newUser()
	return nil
}

func (self *Server) ListUsers() ([]string, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	ret := make([]string, 0, len(self.users))
	for user := range self.users {
		ret = append(ret, user)
	}

	sort.Strings(ret)

	return ret, nil
}

func (self *Server) IsFollowing(who, whom string) (bool, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	uwho, e := self.findUser(who)
	if e != nil {
		return false, e
	}

	_, e = self.findUser(whom)
	if e != nil {
		return false, e
	}

	return uwho.isFollowing(whom), nil
}

func (self *Server) Follow(who, whom string) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	uwho, e := self.findUser(who)
	if e != nil {
		return e
	}

	uwhom, e := self.findUser(whom)
	if e != nil {
		return e
	}

	if uwho.isFollowing(whom) {
		return fmt.Errorf("user %q is already following %q", who, whom)
	}

	uwho.follow(whom, uwhom)
	uwhom.addFollower(who, uwho)
	return nil
}

func (self *Server) Unfollow(who, whom string) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	uwho, e := self.findUser(who)
	if e != nil {
		return e
	}

	uwhom, e := self.findUser(whom)
	if e != nil {
		return e
	}

	if !uwho.isFollowing(whom) {
		return fmt.Errorf("user %q is not following %q", who, whom)
	}

	uwho.unfollow(whom)
	uwhom.removeFollower(who)
	return nil
}

func (self *Server) Following(who string) ([]string, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	uwho, e := self.findUser(who)
	if e != nil {
		return nil, e
	}

	ret := uwho.listFollowing(who)
	return ret, nil
}

func (self *Server) Post(user, at, post string, t time.Time) error {
	if len(post) > trib.MaxTribLen {
		return fmt.Errorf("trib too long")
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	u, e := self.findUser(user)
	if e != nil {
		return e
	}

	if at != "" {
		_, e = self.findUser(at)
		if e != nil {
			return e
		}
	}

	u.post(user, post, self.seq, t)
	self.seq++

	return nil
}

func (self *Server) Home(user string) ([]*trib.Trib, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	u, e := self.findUser(user)
	if e != nil {
		return nil, e
	}

	return u.listHome(), nil
}

func (self *Server) Tribs(user string) ([]*trib.Trib, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	u, e := self.findUser(user)
	if e != nil {
		return nil, e
	}

	return u.listTribs(), nil
}

package ref

import (
	"fmt"
	"sync"
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

func (self *Server) PostTrib(user, post string) error {
	if len(post) > trib.MaxTribLen {
		return fmt.Errorf("trib too long")
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	u, e := self.findUser(user)
	if e != nil {
		return e
	}

	u.post(user, post, self.seq)
	self.seq++

	return nil
}

func (self *Server) FollowedTribs(user string, off, count int) ([]*trib.Trib, error) {
	if off < 0 {
		return nil, fmt.Errorf("negative offset is invalid")
	}
	if count < 0 {
		return nil, fmt.Errorf("negative count is invalid")
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	u, e := self.findUser(user)
	if e != nil {
		return nil, e
	}

	end := off + count
	n := u.ntrib()
	if end >= n {
		end = n
	}
	if off >= n {
		off = n
	}
	return u.list(off, end), nil
}

func (self *Server) CountTribs(user string) (int, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	u, e := self.findUser(user)
	if e != nil {
		return 0, e
	}
	return u.ntrib(), nil
}

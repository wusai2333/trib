package ref

import (
	"time"
	"trib"
)

type user struct {
	following map[string]*user
	followers map[string]*user
	tribs     []*trib.Trib
	timeline  []*trib.Trib
}

func newUser() *user {
	return &user{
		make(map[string]*user),
		make(map[string]*user),
		make([]*trib.Trib, 0, 1024),
		make([]*trib.Trib, 0, 4096),
	}
}

func (self *user) isFollowing(whom string) bool {
	_, found := self.following[whom]
	return found
}

func (self *user) rebuildTimeline() {
	self.timeline = make([]*trib.Trib, 0, 4096)
	for _, user := range self.following {
		self.timeline = append(self.timeline, user.tribs...)
	}
}

func (self *user) follow(whom string, u *user) {
	self.following[whom] = u
	self.rebuildTimeline()
}

func (self *user) unfollow(whom string) {
	delete(self.following, whom)
	self.rebuildTimeline()
}

func (self *user) addFollower(who string, u *user) {
	self.followers[who] = u
}

func (self *user) removeFollower(who string) {
	delete(self.followers, who)
}

func (self *user) post(who string, msg string) {
	t := &trib.Trib{
		User:    who,
		Message: msg,
		Time:    time.Now(),
	}

	// add to my own tribs
	self.tribs = append(self.tribs, t)

	// and it into the timeline of my followers
	for _, user := range self.followers {
		user.timeline = append(user.timeline, t)
	}
}

func (self *user) list(from, to int) []*trib.Trib {
	return self.timeline[from:to]
}

func (self *user) ntrib() int {
	return len(self.timeline)
}

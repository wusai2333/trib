package trib

import (
	"time"
)

const (
	MaxUsernameLen = 15
	MaxTribLen     = 140
)

type Trib struct {
	User    string    // who posted this trib
	Message string    // the content of the trib
	Time    time.Time // the timestamp
}

type Server interface {
	// Creates a user
	SignUp(user string) error

	// Follow someone's timeline
	Follow(who, whom string) error

	// Unfollow
	Unfollow(who, whom string) error // unfollow someone

	// Post a trib
	PostTrib(user, post string) error

	// List the trib of someone's following users
	FollowedTribs(user string, offset, count int) ([]*Trib, error)
}

type Storage interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Append(key, value string) error
	Delete(key string) error
}

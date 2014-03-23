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

	// List all registered users
	ListUsers() ([]string, error)

	// Post a trib
	PostTrib(user, post string) error

	// List the tribs that a particular user posted
	Tribs(user string, offset, count int) ([]*Trib, error)

	// Count of tribs a particular user posted
	CountTribs(user string) (int, error)

	// Returns true if "who" is following "whom"
	IsFollowing(who, whom string) (bool, error)

	// Follow someone's timeline
	Follow(who, whom string) error

	// Unfollow
	Unfollow(who, whom string) error // unfollow someone

	// List the trib of someone's following users
	Home(user string, offset, count int) ([]*Trib, error)

	// Count of tribs for home
	CountHome(user string) (int, error)
}

type Storage interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Append(key, value string) error
	Delete(key string) error
}

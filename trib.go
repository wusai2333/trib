// Package trib defines basic interfaces and constants
// for Tribbler service implementation.
package trib

import (
	"time"
)

const (
	MaxUsernameLen = 15
	MaxTribLen     = 140
	MaxTribFetch   = 100
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
	Post(who, post string, when time.Time) error

	// List the tribs that a particular user posted
	// The result should be sorted in alphabetical order
	Tribs(user string) ([]*Trib, error)

	// Follow someone's timeline
	Follow(who, whom string) error

	// Unfollow
	Unfollow(who, whom string) error // unfollow someone

	// Returns true when who following whom
	IsFollowing(who, whom string) (bool, error)

	// Returns the list of following users
	Following(who string) ([]string, error)

	// List the trib of someone's following users
	Home(user string) ([]*Trib, error)
}

func IsValidUsername(s string) bool {
	if s == "" {
		return false
	}

	if len(s) > MaxUsernameLen {
		return false
	}

	for i, r := range s {
		if r >= 'a' && r <= 'z' {
			continue
		}

		if i > 0 && r >= '0' && r <= '9' {
			continue
		}

		return false
	}

	return true
}

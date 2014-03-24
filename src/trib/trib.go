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
	PostTrib(who, atWhom, post string, when time.Time) error

	// List the tribs that a particular user posted
	Tribs(user string, offset, count int) ([]*Trib, error)

	// Count of tribs a particular user posted
	CountTribs(user string) (int, error)

	// Follow someone's timeline
	Follow(who, whom string) error

	// Unfollow
	Unfollow(who, whom string) error // unfollow someone

	// Returns true if "who" is following "whom"
	IsFollowing(who, whom string) (bool, error)

	// List the trib of someone's following users
	Home(user string, offset, count int) ([]*Trib, error)

	// Count of tribs for home
	CountHome(user string) (int, error)
}

type KeyValue struct {
	Key   string
	Value string
}

func KV(k, v string) *KeyValue { return &KeyValue{k, v} }

type Storage interface {
	Get(key string, value *string) error
	Set(kv *KeyValue, succ *bool) error
	Append(kv *KeyValue, succ *bool) error
	Delete(key string, succ *bool) error
}

func IsValidUsername(s string) bool {
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

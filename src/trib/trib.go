package trib

import (
	"time"
)

const (
	MaxUsernameLen = 15
	MaxTribLen     = 140
)

type Trib struct {
	User    string
	Message string
	Time    time.Time
}

type Server interface {
	SignUp(user string) error
	Follow(who, whom string) error
	Unfollow(who, whom string) error
	PostTrib(user, post string) error
	ListTribs(user string, offset, count int) ([]*Trib, error)
}

type Storage interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Append(key, value string) error
	Delete(key string) error
}

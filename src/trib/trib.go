package trib

import (
	"time"
)

type Trib struct {
	User    string
	Message string
	Time    time.Time
}

type Server interface {
	Register(user string) error
	Subscribe(who, whom string) error
	Unsubscribe(who, whom string) error
	Post(user, post string) error
	List(user string, offset, count int) ([]*Trib, error)
}

type Storage interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Append(key, value string) error
	Delete(key string) error
}

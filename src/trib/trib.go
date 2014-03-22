package trib

import (
	"time"
)

type Server interface {
	Register(user string) error
	Subscribe(who, whom string) error
	Unsubscribe(who, whom string) error
	Post(user, message string) error
	List(user string, offset, count int) ([]*Tribble, error)
}

type Tribble struct {
	Id      uint64
	User    string
	Message string
	Time    time.Time
}

type Storage interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Append(key, value string) error
	Delete(key string) error
}

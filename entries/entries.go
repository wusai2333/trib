package entries

import (
	. "trib"

	"triblab"
)

// Makes a front end that talks to one single backend
// Used in Lab1
func MakeFrontSingle(back string) Server {
	front := &Front{[]string{back}}
	return triblab.NewFront(front)
}

// Serve as a single backend.
// Listen on addr, using s as underlying storage.
func ServeBackSingle(addr string, s Storage, ready chan<- bool) error {
	back := &Back{
		Addr:  addr,
		Store: s,
		Ready: ready,
	}

	return triblab.ServeBack(back)
}

package entries

import (
	. "trib"

	"triblab"
)

// Makes a front end that talks to one single backend
// Used in Lab1
func MakeFrontSingle(back string) Server {
	front := &Front{[]BackGroup{
		BackGroup([]string{back}),
	}}
	return triblab.MakeFront(front)
}

// Makes a front end with multiple backends
// Used in Lab2
func MakeFrontMulti(backs []string) Server {
	front := new(Front)
	front.Backs = make([]BackGroup, len(backs))
	for i, b := range backs {
		front.Backs[i] = BackGroup([]string{b})
	}

	return triblab.MakeFront(front)
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

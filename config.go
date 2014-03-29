package trib

// Backend config
type Back struct {
	Addr  string      // listen address
	Store Storage     // the underlying storage it should use
	Ready chan<- bool // send a value when server is ready
}

// Peering parameters, used in Lab2/3
type Peering struct {
	// Non zero incarnation identifier, 0 if have no peers
	Id int

	// The address of this backend, empty string if have no peers
	You string

	// The addresses of peers including you, nil if have no peers
	Peers []string
}

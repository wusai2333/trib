package trib

// Frontend config
type Front struct {
	Backs []string // List of backend addresses
}

// Backend config
type Back struct {
	Addr  string      // listen address
	Store Storage     // the underlying storage it should use
	Ready chan<- bool // send a value when server is ready

	// The following are peering parameters, used in Lab2/3 only

	// The address of this backend, empty string if have no peers
	You string

	// The addresses of peers including you, nil if have no peers
	Peers []string

	// Non zero incarnation identifier, 0 if have no peers
	Id int
}

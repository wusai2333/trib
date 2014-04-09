package trib

// Backend config
type BackConfig struct {
	Addr  string      // listen address
	Store Storage     // the underlying storage it should use
	Ready chan<- bool // send a value when server is ready
}

type KeeperConfig struct {
	// The addresses of peers including the address of this back-end
	Addrs []string

	// The index of this back-end
	This int

	// Non zero incarnation identifier
	Id int64
}

func (c *KeeperConfig) Addr() string {
	return c.Addrs[c.This]
}

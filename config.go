package trib

// Backend config
type BackConfig struct {
	Addr  string      // listen address
	Store Storage     // the underlying storage it should use
	Ready chan<- bool // send a value when server is ready

	Peer *PeerConfig // only used in Lab2 and Lab3
}

type PeerConfig struct {
	// The addresses of peers including the address of this back-end
	Addrs []string

	// The index of this back-end
	This int

	// Non zero incarnation identifier
	Id int64
}

func (p *PeerConfig) Addr() string {
	return p.Addrs[p.This]
}

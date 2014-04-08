package trib

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type RC struct {
	Backs []*BackAddr
}

type BackAddr struct {
	Serve string
	Peer  string
}

func (self *RC) ServeAddrs() []string {
	ret := make([]string, 0, len(self.Backs))

	for _, b := range self.Backs {
		ret = append(ret, b.Serve)
	}

	return ret
}

func (self *RC) PeerAddrs() []string {
	ret := make([]string, 0, len(self.Backs))

	for _, b := range self.Backs {
		ret = append(ret, b.Peer)
	}

	return ret
}

func (self *RC) BackCount() int {
	return len(self.Backs)
}

func (self *RC) BackConfig(i int, s Storage) *BackConfig {
	ret := new(BackConfig)
	back := self.Backs[i]
	ret.Addr = back.Serve
	ret.Store = s
	ret.Ready = make(chan bool, 1)

	ret.Peer = new(PeerConfig)
	ret.Peer.Addrs = self.PeerAddrs()
	ret.Peer.This = i
	ret.Peer.Id = time.Now().UnixNano()

	return ret
}

func LoadRC(p string) (*RC, error) {
	fin, e := os.Open(p)
	if e != nil {
		return nil, e
	}
	defer fin.Close()

	ret := new(RC)
	e = json.NewDecoder(fin).Decode(ret)
	if e != nil {
		return nil, e
	}

	return ret, nil
}

func (self *RC) marshal() []byte {
	b, e := json.MarshalIndent(self, "", "    ")
	if e != nil {
		panic(e)
	}

	return b
}

func (self *RC) Save(p string) error {
	b := self.marshal()

	fout, e := os.Create(p)
	if e != nil {
		return e
	}

	_, e = fout.Write(b)
	if e != nil {
		return e
	}

	_, e = fmt.Fprintln(fout)
	if e != nil {
		return e
	}

	return fout.Close()
}

func (self *RC) String() string {
	b := self.marshal()
	return string(b)
}

package ready

import (
	"net/rpc"
)

func Notify(addr, s string) error {
	c, e := rpc.DialHTTP("tcp", addr)
	if e != nil {
		return e
	}

	var b bool
	e = c.Call("Ready.Ready", s, &b)
	if e != nil {
		return e
	}

	return c.Close()
}

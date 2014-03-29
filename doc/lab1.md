## Lab1

Welcome to Lab1. In this lab, you need to do several things.

1. Implement a key-value storage server that fits `trib.Store` interface
   and takes the Http RPC requests from the network. 
2. Implement a key-value storage client that fits `trib.Store` interface
   but calls a remote RPC server.
3. Implement a stateless Tribbler front-end that fits `trib.Server` interface
   but calls a remote RPC server.

In specific, you need to implement the 3 functions that are currently defined
in `trib/entries.go` file: `ServeBack()`, `NewClient()` and `NewFront()`. They
are now all filled with a one-line todo place holder.

```
func ServeBack(b *trib.Back) error
```
This function creates an instance of a back-end server based on configuration
`b *trib.Back`. Structure `trib.Back` is defined in `trib/config.go` file.
In the struct type, it has several fields:

- `Addr` is the address the server should listen on, in
the form of `<host>:<port>`. Go language uses this address in its [standard
`net` package] (http://golang.org/pkg/net), so you should be able to use it
directly.  
- `Store` is the storage device you will use for storing the key-value
pair. 
- `Ready` is a channel for notifying the other parts in the program that the 
server is ready to accept RPC calls from the network. The value that you
send into the tunnel does not matter.

```
func NewClient(addr string) trib.Stroage
```

This function takes the addr as a TCP address in the form of `<host>:<port>`,
and will use that as the server address. It returns an implementation of
`trib.Storage`.

```
func NewFront(f *trib.Front) trib.Server


## RPC

Go language has its own [`net/rpc` package](http://golang.org/pkg/net/rpc),
and we will use that. Note that the `trib.Store` interface is already
in its "RPC friendly" form.



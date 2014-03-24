# tribble

- trib                  // defines the interfaces
- trib/cmd/client       // a sample client that performs http calls
- trib/cmd/trib-front   // the webserver, binds with a particular front end implementation, with rate limiter in it
- trib/cmd/trib-back    // the simple storage backend program, binds uses trib/back or trib/dupback
- trib/ref              // reference implementation, defines the service, working but not scalable in any sense
- trib/store       // storage service lib, with rate limiter and error generator in it
- triblab/lab1     // lab1: the stateless frontend logic package
- triblab/lab2     // lab2: rpc bridge
- triblab/lab3     // lab3: duplicated backend storage
- tribtest/lab1	   // black box test cases
- tribtest/lab2
- tribtest/lab3 

type Tribble struct {
	Id uint64
    User string
    Message string
    Time time.Time
}

type Storage inteface {
    Get(key string) (string, error)
    Set(key, value string) error
    Append(key, value string) error
    Delete(key string) error
}

// message length: 140 runes
// username length: 15 runes
// timeline fetch max: 100 tribs

type Server interface {
    Register(user string) error
    Subscribe(who string, whom string) error
    Unsubscribe(who string, whom string) error
    Post(user, message string) error
    List(user string, offset, count int) ([]*Tribble, error)
}

ref.NewServer() Server

// lab1: decouple the storage, using a key-value pair storage
// attention: make the logic stateless, and robust to failure, have some retries
lab.NewFront(backend string) (Server, error)
lab.ServeBack(addr string, store Storage) error

type Frontend struct {
	Backends []string
}

lab.NewFront(backends []string) (Server, error) // connects
lab.ServeBack(addr string, s Storage)
lab.Backend.Serve() error // RPC storage server

// lab2: rpc to the storage, stateless front-end
// attention: handle errors
lab.NewServer(backends []string) (Server, error) // connects 

// lab3:
type Backend struct {
	Addr string
	Store Storage

    Peers []string
    You string
    Id int
}
Backend.Serve() error

// lab1: implement the service logic
// lab2: make the backend interface rpc
// lab3: make the backend a duplicated service

For the vector clock 

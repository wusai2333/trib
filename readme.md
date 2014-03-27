# tribble

- trib                  // defines the interfaces
- trib/cmd/trib-front   // the webserver, binds with a particular front end implementation, with rate limiter in it
- trib/cmd/trib-back    // the simple storage backend program, binds uses trib/back or trib/dupback
- trib/ref              // reference implementation, defines the service, working but not scalable in any sense
- trib/store            // storage service lib, with rate limiter and error generator in it
- trib/entries          // how triblab will be called for different labs
- triblab               // the labs
- tribtest/lab1	        // black box test cases

// lab1: implement the service logic
// lab2: make the backend interface rpc
// lab3: make the backend a duplicated service

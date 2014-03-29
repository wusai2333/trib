## Lab1

Welcome to Lab1. The goal of this lab is to split the logic into a stateless
scalable front-end and a key-value pair backend. In particular, you need
to:

1. Implement a key-value storage server type that fits `trib.Store` interface
   and takes http RPC requests from the network.
2. Implement a key-value storage client type that fits `trib.Store` interface
   that calls a remote RPC key-value pair server.
3. Implement a stateless Tribbler front-end type that fits `trib.Server` interface
   that calls a remote RPC key-value pair back-end server.

More specifically, you need to implement three entry functions that are
defined in `triblab/entries.go` file: `ServeBack()`, `NewClient()` and
`NewFront()`. They are now filled with `panic("todo")`.

## Tribble

A Tribble is a structure that has 4 fields:

```
type Trib struct {
    User    string    // who posted this trib
    Message string    // the content of the trib
    Time    time.Time // the timestamp
    Clock   uint64    // a logical clock, not used in lab1
}
```

`Time` is what the front-end claims when
this tribble is created. However, to sort
tribbles in a globally consistent and reasonable order,
Tribbler service maintains a distributed logical `Clock` in `uint64`.

When sorting tribbles into a timeline, you should follow this field priroty:

1. `Clock`
2. `Time`
3. `User`
4. `Message`

We call this *Tribble Order*.

## Tribbler Service Interface

The Tribbler service logic is all defined in `trib.Server` interface
(in `trib/trib.go` file).

*** 

```
SignUp(user string) error
```

Creates a new user. After a user is created, it will always be there.
A user name must be no longer than `trib.MaxUsernameLen = 15`
characters but non-empty, must start with a lower-case letter,
and can only contain lower-case letters or numbers.

There is a helper function called `trib.IsValidUsername(string)`
which you can use to check if a username is valid.

Returns error when the user already exists. Concurrent sign-ups might
both succeed.

***

```
ListUsers() ([]string, error)
```

List `trib.MinListUser = 20` registered users.  When there are less than 20
users that have ever signed up, list all of them. This is just for showing some
users on the front page. This is not for listing all the users that signed up,
because that would be too expensive.

***

```
Post(who, post string, clock uint64) error
```

Post a tribble. `clock` is the maximum clock value this user
has ever seen so far by reading tribbles
(via `Home()` and `Tribs()`).
It returns error when the user does not exist or the post
is too long (longer than `trib.MaxTribLen = 140`).

***

```
Tribs(user string) ([]*Trib, error)
```

List the recent `trib.MaxTribFetch = 100` tribbles that a user
posted. Tribbles needs to be sorted in Tribble Order. Also,
it should make sure that the order is the same order
that the user posted the tribbles.

***

```
Follow(who, whom string) error
Unfollow(who, whom string) error
IsFollowing(who, whom string) (bool, error)
Following(who string) ([]string, error)
```

Functions for actions like follow/unfollow, check following
and listing all following users. A user can never
follow or unfollow himself. When calling with `who` equals
to `whom`, the functions return error. When the user
does not exist, the functions return error.

***

```
Home(user string) ([]*Trib, error)
```

List the recent `trib.MaxTribFetch` = 100 tribbles that are
posted on the user's following timeline in Tribble Order.
In addition, the order should always satisfy:

1. If a tribble A is posted 3 seconds after a tribble B is
posted, A always shows after B.
2. If a tribble A is posted after a user sees tribble B,
A always shows after B.

It returns error when the user does not exist.

***

In addition to normal errors, it might also return IO errors
if the implementation needs to communicate to a remote part.
Returning a nil error means the the call is successfully
executed; returning a error that is not nil means that
the call might be executed or not.

## Key-value Pair Service Interface

Data structure and interfaces for the key-value pair service is defined
in `trib/kv.go` file. The main interface is `trib.Storage` interface,
which consists of three parts.

First is the key-value pair part.

```
// Key-value pair interfaces
// Default value for all keys is empty string
type KeyString interface {
	// Gets a value. Empty string by default.
	Get(key string, value *string) error

	// Set kv.Key to kv.Value. Set succ to true when no error.
	Set(kv *KeyValue, succ *bool) error

	// List all the keys of non-empty pairs where the key matches
	// the given pattern.
	Keys(p *Pattern, list *List) error
}
```

`Pattern` is a prefix-suffix tuple. It has a `Match(string)` function
that returns true when the string matches the pattern.

Second is the key-string pair part.

```
// Key-list interfaces.
// Default value for all lists is an empty list.
// After the call, list.L should never by nil.
type KeyList interface {
	// Get the list.
	ListGet(key string, list *List) error

	// Append a string to the list, succ will always set to true.
	ListAppend(kv *KeyValue, succ *bool) error

	// Removes all elements that equals to kv.Value in list kv.Key
	// n is set to the number of elements removed.
	ListRemove(kv *KeyValue, n *int) error

	// List all the keys of non-empty lists, where the key matches
	// the given pattern.
	ListKeys(p *Pattern, list *List) error
}
```

And finally we put it together with an auto-incrementing clock service:

```
type Storage interface {
	// Returns an auto-incrementing clock, the returned value
	// will be no smaller than atLeast, and it will be
	// strictly larger than the value returned last time,
	// unless it was math.MaxUint64.
	Clock(atLeast uint64, ret *uint64) error

	KeyString
	KeyList
}
```

Note that the function signature of these methods are all RPC friendly.  Also,
under the defintion of the execution logic, all the methods will always return
nil error. This means all errors you see from this interface will be
communication errors. 

## Entry Functions

These are the 3 entry functions you need to implement.

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
- `Store` is the storage device you will use for storing stuff. In fact,
You should not store persistent data anywhere else.
- `Ready` is a channel for notifying the other parts in the program that the
server is ready to accept RPC calls from the network. The value that you
send into the tunnel does not matter.

This function should be a blocking call. It does not return until an error
(like the network is shutdown) occurred.

```
func NewClient(addr string) trib.Stroage
```

This function takes the addr as a TCP address in the form of `<host>:<port>`,
and will use that as the server address. It returns an implementation of
`trib.Storage`. You can assume `addr` will always be an valid address.

```
func NewFront(backs []string) trib.Server
```

This function takes the addresses of the backends, and returns an implementation
of `trib.Server`. The returned instance then will serve as an service front-end
that takes Tribbler service function calls, and translates them into key-value
pair RPC calls. This front-end should be stateless, thread safe, and ready
to be killed at any time. This means that at any time during its execution,
the back-end key-value pair storage always stays in an consistent. Also, note
that one front-end might be taking multiple

## RPC

Go language comes with its own
[`net/rpc`](http://golang.org/pkg/net/rpc) package
in its standard library, and we will just use that.
Note that the `trib.Store` interface is already in its "RPC friendly" form.

Your RPC needs to use the default encoding `encoding/gob`, listen on the given
address, and serve as an http RPC server.

## Testing

Both the `trib` and `triblab` repository comes with a makefile with some handy
command lines, and also some basic testing code.

Under `trib` directory, if you type `make test`, you should see 
that the tests runs and all passed.

Under `triblab` directory, if you type `make test` however, you would see
the test fails with a todo panic.

Your first attempt should be implement the logic and try to pass those test
cases. If you pass those, you should be fairly confident that you can
get at least 30% of the credits for Lab1 (unless you are cheating in some way).

However, the test that comes with the repository is very basic and simple.
Though you don't have to, you should really write more test cases to make sure
your implementation matches the specification.

## Playing with It

To run your own implementation, you could use the `trib-front` and `trib-back`
launcher.

First make sure you code compiles.

Then run the back-end server.

```
$ trib-back
```

And you should see an address printing out, say it is `localhost:37021`.
Note that you can also specify your own address via command line.
The default is `localhost:rand`.

Next for the front-end part.

```
$ trib-front -init -lab -back=localhost:37021
```

For the `-back` flag, please use the backend address that you just got
from running `trib-back`.

`-init` will populate the service with some sample data.
`-lab` tells the front-end to connect to a back-end rather than running
with the default reference implementation.

Now you can open your browser, connect to the front-end machine
and play with your own implementation.

Note that, when you completes Lab1, it should be perfectly fine to have
multiple front-ends that connects to a single back-end.

## Performance Requirement

When running on the lab machines, the tribbler front-end service should return
every function call within 1 second.

## Common Mistakes

Here are some common mistakes that a lazy and quick
but incorrect implementation might do:

- **Read-modify-write**. For example, a tribbler might read a counter
  from the key-value store, increase it by one
  and then write it back (at the same key). 
  This will introduce racing condition among the front-ends.
- **Not handling errors**. A tribbler service call might require several
  RPC calls to the backend. It is important to properly handle *any* error
  returned by these calls.
- **Sorting by the timestamps first**. Again, Tribble Order means that
  the logic clock is the first field to consider. 
- **Use the clock argument from the front-end for the clock
  field of a new Tribble**. Well, technically, you can do that in your
  code internally as long as you can satisfy the ordering requirements
  speficied for `Home()` and `Tribs()` (you might find it very hard).
  Nonetheless, intuitively, the clock argument tells the *oldest* tribble a
  user have seen (which might be 0 if the user has not seen any tribble yet), 
  hence the new posted tribble seems to better have a clock value that is
  larger than the argument.
- **Generate the clock from the timestamp**. While 64-bit can cover a very
  wide time range even in the unit of nanoseconds, you should keep in mind
  that the front-ends are running on different servers with arbitrary physical
  time differences, so it is not wise to generate the logical *clock* from the
  physical *time*.
- **Not handling old tribbles**. Note that only the most recent 100 tribbles
  of a user matters. Not handling old tribbles might lead to worse and worse
  performance over time and eventually break the performance promise.

## Turning In

First, make sure every piece of your code is commited into the repository in
`triblab`. Then just type `make turnin` under the root of the repository.
It will generate a turnin.zip that contains everything in you repo, and
copy that out.

## Happy Lab1!

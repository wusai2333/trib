## Lab1

Welcome to Lab1. The goal of this lab is to split the logic into a stateless
scalable front-end and a key-value pair backend. In particular, you need
to:

1. Implement a key-value storage server type that fits `trib.Store` interface
   and takes the Http RPC requests from the network.
2. Implement a key-value storage client type that fits `trib.Store` interface
   but calls a remote RPC server.
3. Implement a stateless Tribbler front-end type that fits `trib.Server` interface
   but calls a remote RPC server.

More specifically, you need to implement the 3 entry functions that are currently defined
in `trib/entries.go` file: `ServeBack()`, `NewClient()` and `NewFront()`. They
are now all filled with a one-line todo place holder.

## Tribble

A Tribble is a type that has 4 fields:

```
type Trib struct {
	User    string    // who posted this trib
	Message string    // the content of the trib
	Time    time.Time // the timestamp
	Clock   uint64    // a logical clock, not used in lab1
}
```

Timestamp is what the front-end claims the
time that this tribble is created. However, for sorting
tribbles in a globally consistent and reasonable order,
Tribbler service maintains a logical clock in `uint64`.

When sorting tribble timelines, one should first
order them by the `Clock` field, and then by the `Time` field,
then by the `User` field, and finally by the `Message` fields.
For most of the cases, `Clock` and `Time` would do the work.
We call this *Tribble Order*.

## Tribbler Service

The Tribbler service logic is all defined in `trib.Server` interface
(in `trib/trib.go` file).

```
SignUp(user string) error
```

Creates a new user. After a user is created, it is always there.
A user name must be no longer than `trib.MaxUsernameLen` = 15
characters but not empty, must start with a lower-case letter,
and can only contain lower-case letters or numbers

There is a helper function called `trib.IsValidUsername(string)`
which you can use to check if a username is valid.

When the user exists, it returns error.

```
ListUsers() ([]string, error)
```

List `trib.MinListUser` = 20 registered users.
When there are less than 20 users that
signed up, list all of them. This is for showing some users on
the front page. This is not for listing all the users that signed
up, because that would be too expensive.

```
Post(who, post string, clock uint64) error
```

Post a tribble. `clock` is the maximum clock value this user
has ever seen so far by reading tribbles
(via `Home()` and `Tribs()`).
It returns error when the user does not exist or the post
is too long (longer than `trib.MaxTribLen` = 140).

```
Tribs(user string) ([]*Trib, error)
```

List the recent `trib.MaxTribFetch` = 100 tribbles that a user
posted. Tribbles needs to be sorted in Tribble Order. Also,
it should make sure that the order is the same order
that the user posted the tribbles.

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

--

In addition to normal errors, it might also return IO errors
if the implementation needs to communicate to a remote part.
Returning a nil error means the the call is successfully
executed; returning a error that is not nil means that
the call might be executed or not.

## Key-value Pair Service



## Entries

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
[`net/rpc` package](http://golang.org/pkg/net/rpc),
in its standard library, and we will just use that.
Note that the `trib.Store` interface is already in its "RPC friendly" form.

Your RPC needs to use the default encoding `encoding/gob`, listen on the given
address, and serve as an Http RPC server.

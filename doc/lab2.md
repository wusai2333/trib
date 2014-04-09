## Lab2

Welcome to Lab2. The goal of this lab is to use the RPC key-value pair
we built in Lab 1 to implement scalable Tribbler front-ends and
back-ends. In particular, you need to implement a stateless Tribbler
front-end type that fits `trib.Server` interface and can perform calls
to a list of remote RPC key-value pair back-end server, which save all
the user information and Tribbles in a distributed fashion.

## Get Your Repo Up-to-date

First we update the `trib` repo:

```
$ cd ~/gopath/src/trib
$ git pull /classes/cse223b/sp14/labs/trib lab2
```

If you have not changed anything in `trib` repo, this should be
painless. However, if you changed stuff, you need to merge the
changes.

Now update the `triblab` repo by merging branch `lab2`. There will be
several changes:

- Some line changes in `makefile`.
- Some added lines in `lab2.go`.
- A new file called `server_test.go`.

If you have not touched those files and have not created a file called
`server_test.go` by yourself, the merge should be painless:

```
$ cd ~/gopath/src/triblab
$ git pull /classes/cse223b/sp14/labs/triblab lab2
```

If you have made changes to those files. Then you need to merge the
changes properly.

## Tribble

A Tribble is a structure type that has 4 fields:

```
type Trib struct {
    User    string    // who posted this trib
    Message string    // the content of the trib
    Time    time.Time // the physical timestamp
    Clock   uint64    // the logical clock
}
```

`Time` is what the front-end claims when this tribble is created, by
reading the front-end's own physical time clock on the machine when
`Post()` in a `trib.Server` is called.  However, to sort tribbles in a
globally consistent and reasonable order, we can not sort the tribbles
only by this timestamp, because different front-ends have different
physical time readings. For sorting, Tribbler service needs to
maintain a distributed logical `Clock` in `uint64`.

When sorting many tribbles into a single timeline, you should sort by
the fields following this priority:

1. `Clock` The logical timestamp.
2. `Time` The physical timestamp.
3. `User` The user id
4. `Message` The message content.

We call this the *Tribble Order*.

## Tribbler Service Interface

The Tribbler service logic is all defined in `trib.Server` interface
(in `trib/trib.go` file). This is how the webpage user interface
interacts with a Tribbler server.

***

```
SignUp(user string) error
```

Creates a new user. After a user is created, it will never disappear
in the system.  

A valid user name must be no longer than `trib.MaxUsernameLen=15`
characters but not empty, must start with a lower-case letter, and can
only contain lower-case letters or numbers.  There is a helper
function called `trib.IsValidUsername(string)` which you can use to
check if a username is valid.

Returns error when the username is invalid or the user already exists.
Concurrent sign-ups might both succeed.

***

```
ListUsers() ([]string, error)
```

Lists at least `trib.MinListUser=20` different registered users. When
there are less than 20 users that have ever signed up, list all of
them. The returned usernames should be sorted in alphabetical order.

This is just for showing some users on the front page; this is not for
listing all the users that have ever signed up, because that would be
too expensive in a scalable system.

***

```
Post(who, post string, clock uint64) error
```

Posts a tribble. `clock` is the maximum clock value this user client
has ever seen so far by reading tribbles (via `Home()` and `Tribs()`).
It returns error when the user does not exist or the post is too long
(longer than `trib.MaxTribLen=140`).

***

```
Tribs(user string) ([]*Trib, error)
```

Lists the recent `trib.MaxTribFetch=100` tribbles that a user posted.
Tribbles needs to be sorted in the Tribble Order. Also, it should make
sure that the order is the same order that the user posted the
tribbles.

***

```
Follow(who, whom string) error
Unfollow(who, whom string) error
IsFollowing(who, whom string) (bool, error)
Following(who string) ([]string, error)
```

These are functions to follow/unfollow, check following and listing
all following users of a user. A user can never follow or unfollow
himself. When calling with `who` equals to `whom`, the functions
return error. When the user does not exist, the functions return
error.

A user can follow at most `trib.MaxFollowing=2000` users. Returns
error when trying to follow more than that.

***

```
Home(user string) ([]*Trib, error)
```

List the recent `trib.MaxTribFetch=100` tribbles that are posted on
the user's following timeline in Tribble Order.  In addition, the
ordering should always satisfy that:

1. If a tribble A is posted after a tribble B is posted, and they are
posted by the same user, A always shows after B.
2. If a tribble A is posted 10 seconds after a tribble B is posted,
even if they are posted by different users, A always shows after B.  
3. If a tribble A is posted after a user client sees tribble B, A
always shows after B.

A is *posted after* B means B calls `Post()` after A's `Post()`
returned.

It returns error when the user does not exist.

***

In addition to normal errors, it might also return IO errors if the
implementation needs to communicate to a remote part.  Returning a nil
error means that the call is successfully executed; returning a
non-nil error means that the call might be succefully executed or not.

## System Architecture

The system architecture looks like this:

![System Arch](./arch.png)

Users use the Tribbler system from their browsers, each user visiting
one state-less front-end at a time (probably distributed via DNA
multiplexing). Upon a service request, the front-end will translate
the request into key-value pair requests and issue these requests to
the back-ends over RPC. The back-ends can talk to each other via
back-end peering channels (which you will implement with your own
design), so that they can sync on their view of the world from
time to time. 

The peering channels should serve for these purposes:

1. Coarse grained time synchronization, so that the clock in the system
   won't be offset too much.
2. Fault detection, so that when a back-end joins or leaves, the
   system can adjust itself.
3. Consistent storage management, so that when a back-end joins or
   leaves, the system won't lose any-data and keeps the service
   in a consistent view.

## Entry Functions

You can find these entry functions in `lab2.go` file under
`triblab` repo:

```
func NewFront(backs []string) trib.Server
```

This function takes the addresses of the backends, and returns an
implementation of `trib.Server`. The returned instance then will serve
as an service front-end that takes Tribbler service function calls,
and translates them into key-value pair RPC calls. This front-end
should be stateless, thread safe, and ready to be killed at any time.
This means that at any time during its execution, the back-end
key-value pair storage always needs to stay in a consistent state.
Also, note that one front-end might be taking multiple
concurrent requests from the Web,
and there might be multiple front-ends talking to the same
back-end, so make sure it handles all the concurrency issues
correctly.

In addition to `NewFront()`, you also need to make changes to
`ServeBacks()` in `lab1.go1`.

```
func ServeBacks(b *trib.BackConfig) error
```

The signature remains unchanged, but we added a new field in
`trib.BackConfig` called `Peer`. Now `trib.BackConfig` looks like
this

```
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
```

Here explains the new fields in `Peer`:

- `Addrs` are the addresses where the back-ends will listen on for
  peering. This is different from `Addr` in `BackConfig`. Theses
  addresses are only for back-ends to communicate with each other,
  where the `Addr` is for a front-end (like a `kv-client`) to connect
  into this back-end.
- `This` indicates the index of this back-end in the `Addrs`, so the
  back-end should listen on `Addrs[This]` for connections from other
  back-ends.
- `Id` is a unique incarnation identifier for this backend. It is not
  particularly useful in Lab2, but will be useful in Lab3 when the
  back-ends need to be fault-tolerant.

You can design your own protocol for communication
among the back-ends.

## Playing with It

Since we might have multiple back-ends running at the same time.
To work with that, we will have a configuration file that
specifies the serving address (which accepts connections from front-ends)
and the peering address (which accepts connections from back-end peers).
The config file by default will use `trib.rc` as its file name.

`trib.rc` is saved in json format, marshalling a `RC` structure type
(defined in `trib/rc.go` file).
We have a utility program called `trib-mkrc`
that can generate a `trib.rc` file
automatically.

Find a directory as your working directory (like `triblab`).

```
$ trib-mkrc -local -n=3
```

This will generate a file called `trib.rc` under the current
directory, and also print the file content to stdout.  `-local` means
that all addresses will be on `localhost`.  `-n=3` means there are in
total 3 back-ends.  If you remove `-local`, then it will generate
back-ends starting from `172.22.14.211` through `172.22.14.220`, which
are the IP address of our lab machines. There can be 10 backends in
maximum.

With this configuration file, we can now launch the
back-ends:

```
$ trib-back
```

This will read and parse the `trib.rc` file, and spawn all the
back-ends which serving addresses are on this host. Since all the
back-ends we generate here are on `localhost`, so all the back-ends
are spawned for this case (in different go routines). You should see
three log lines showing that three back-ends just started.

Next for the front-end part:

```
$ trib-front -init -addr=:rand -lab
```

You have used this utility before. The only new thing here is the
`-lab` flag, which tells it to read the `trib.rc` file and use our lab
implementation. This will start a stateless front-end (which you
implemented in this lab) that will connect to the back-ends service
addresses specified in `trib.rc`.

Again `-init` will populate the service with some sample data.

Now you can open your browser, connect to the front-end machine and
play with your own implementation.

If you want to use some other config file, use the `-rc` flag.
It is supported in all `trib-*` utilities.

Note that, when you completes this lab, it should be perfectly fine to
have multiple front-ends that connects to the set of back-ends.
Also both the front-ends and the back-ends should be scalable.

## Assumptions

These are some unreal assumptions you can have for Lab2.

- No network communication error will happen.
- Once a back-end starts, it will remain online forever.
- The `trib.Storage` used in the backend will return every `Clock()`
  call in less than 1 second.
- In the `trib.Storage` used in the backend, each key visiting
  (checking if the key exist, locating its corresponding value, or as
  a process of iterating all keys) will take less than 1 millisecond.
  Read and write 1MB of data on the value part (in list or string)
  will take less than 1 millisecond.  Note that `Keys()` and
  `ListKeys()` might take longer time to complete because it needs to
  scan over all the keys.
- All back-end servers will run on the lab machines.
- Although a front-end can be killed at any-time, the killing only
  happens very occasionally.

Note that some of them won't stay in Lab3, so
try not to rely on the assumptions too much.

## Requirements

- The back-ends should be able to start one-by-one in arbitrary order.
  without any error.
- When the service function call has valid arguments, the function
  call should not return any error.
- The front-ends part should be stateless and hence ready to be killed
  at anytime.
- When the back-end is the system throughput bottleneck, adding more
  back-ends should increase the number of service function calls the
  system can serve per second.
- When running on the lab machines, with more than 5 back-ends
  supporting, each Tribbler service call should return in 3 seconds.
- Each back-ends should maintain the same general key-value pair
  semantics as they were in Lab1. As a result, all test cases that
  pass for Lab1 should also pass for Lab2. This means that, the
  back-ends do not need to understand anything about the front-ends
  (like how the keys will be structured and organized, or how to parse
  the values).

## Building Hints

While you are free to build the front-ends and the back-ends in
your own way, here are some suggested hints:

- For each service call in the front-end, if it updates anything in
  the back-end storage, use only one write-RPC call. This will make
  sure it the call either succeed or fail.
- Hash the tribbles and other information into all the back-ends based
  on username. You may find the package `hash/fnv` helpful for
  hashing.
- Synchronize the logical clocks among all the back-ends every second.
  (This will also serve as a heart-beat signal, which will be useful
  for implementing Lab3.) However, you should not try to synchronize
  the clocks for every post, because that will be not scalable.
- Do some garbage collection when one user have too many tribbles
  saved in the storage.
- Keep a cache for the ListUsers() call when the users are many.

## Common Mistakes

Here are some common mistakes that a lazy and quick
but incorrect implementation might do:

- **Read-modify-write**. For example, a tribbler might read a counter
  from the key-value store, increase it by one and then write it back
  (at the same key).  This will introduce racing condition among the
  front-ends.
- **Not handling errors**. A tribbler service call might require
  several RPC calls to the backend. It is important to properly handle
  *any* error returned by these calls. It is okay to tell the webpage
  that an error occurred. However, it is often not a good idea to leave
  the back-end in inconsistent state.
- **Sorting by the timestamps first**. Again, the Tribble Order means
  that the logic clock is the first field to consider on sorting.
- **Misuse the clock argument in Post()**. For example, you
  might directly use that argument as the new post's clock field.
  Technically, you can do that in your code internally as long as you
  can satisfy the ordering requirements specified for `Home()` and
  `Tribs()` (you might find it very hard).  Nonetheless, intuitively,
  the clock argument tells the *oldest* tribble a user have seen
  (which might be 0 if the user has not seen any tribble yet), hence
  the new posted tribble seems to better have a clock value that is
  larger than the argument.
- **Generate the clock from the timestamp**. While 64-bit can cover a
  very wide time range even in the unit of nanoseconds, you should
  keep in mind that the front-ends are running on different servers
  with arbitrary physical time differences, so it is not wise to
  generate the logical *clock* from the physical *time*.
- **Not handling old tribbles**. Note that only the most recent 100
  tribbles of a user matter. Not handling old tribbles might lead to
  worse and worse performance over time and eventually break the
  performance promise.

## Turning In

First, make sure that you have committed every piece
of your code into the repository `triblab`. Then just
type `make turnin` under the root of the repository.
It will generate a `turnin.zip` that contains everything
in your git repository, and will then copy the zip file to
a place where only the lab instructors can read.

## Happy Lab2

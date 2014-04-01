# Lab2

Welcome to Lab2. The goal of this labe is to ...
The goal of this lab is to split the logic into
stateless scalable front-ends and a key-value pair backend. In
particular, you need to:

3. Implement a stateless Tribbler front-end type that fits
`trib.Server` interface that calls a remote RPC key-value pair
back-end server.

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

`Time` is what the front-end claims when this tribble is created,
by reading the front-end's own physical time clock on the machine when
`Post()` is called. However, to sort tribbles in a globally 
consistent and reasonable order, this timestamp cannot be used
because different front-ends have different physical time readings.
Tribbler service maintains a distributed logical `Clock` in `uint64`
for sorting.

When sorting many tribbles into a single timeline, you should sort by
the fields following this priority:

1. `Clock` The logical timestamp.
2. `Time` The physical timestamp.
3. `User` The user id
4. `Message` The message content.

We call this the *Tribble Order*.

## Tribbler Service Interface

The Tribbler service logic is all defined in `trib.Server` interface
(in `trib/trib.go` file).

*** 

```
SignUp(user string) error
```

Creates a new user. After a user is created, it will always exist.  A
user name must be no longer than `trib.MaxUsernameLen=15` characters
but not empty, must start with a lower-case letter, and can only
contain lower-case letters or numbers.

There is a helper function called `trib.IsValidUsername(string)` which
you can use to check if a username is valid.

Returns error when the username is invalid or the user already exists.
Concurrent sign-ups might both succeed.

***

```
ListUsers() ([]string, error)
```

Lists at least `trib.MinListUser = 20` different registered users.
When there are less than 20 users that have ever signed up, list all
of them. The returned usernames should be sorted in alphabetical.

This is for showing some users on the front page.  This is not for
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
Tribbles needs to be sorted in Tribble Order. Also, it should make
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
order should always satisfy:

1. If a tribble A is posted after a tribble B is posted, and they are
posted by the same user, A always shows after B.
2. If a tribble A is posted 5 seconds after a tribble B is posted,
even if they are posted by two different users, A always shows after
B.  
3. If a tribble A is posted after a user client sees tribble B, A
always shows after B.

A is *posted after* B means B calls `Post()` after A's `Post()`
returned.

It returns error when the user does not exist.

***

In addition to normal errors, it might also return IO errors if the
implementation needs to communicate to a remote part.  Returning a nil
error means that the call is successfully executed; returning a
non-nill error means that the call might be executed or not.

## Entry Functions

```
func NewFront(backs []string) trib.Server
```

This function takes the addresses of the backends, and returns an
implementation of `trib.Server`. The returned instance then will serve
as an service front-end that takes Tribbler service function calls,
and translates them into key-value pair RPC calls. This front-end
should be stateless, thread safe, and ready to be killed at any time.
This means that at any time during its execution, the back-end
key-value pair storage always stays in a consistent state. Also, note
that one front-end might be taking multiple concurrent requests from
the Web, and there might be multiple front-ends talking to the same
back-end, so make sure it handles all the concurrency issues
correctly.

In Lab1, `backs` will always contain only one address.

## Playing with It

Next for the front-end part. Run on the same machine (since it is
listening on `localhost`).

```
$ trib-front -init -lab -back=localhost:37021
```

For the `-back` flag, please use the backend address that you just got
from running `trib-back`.

`-init` will populate the service with some sample data.  `-lab` tells
the front-end to connect to a back-end rather than running with the
default reference implementation.

Now you can open your browser, connect to the front-end machine and
play with your own implementation.

Note that, when you completes Lab1, it should be perfectly fine to
have multiple front-ends that connects to a single back-end.

## Common Mistakes

Here are some common mistakes that a lazy and quick
but incorrect implementation might do:

- **Read-modify-write**. For example, a tribbler might read a counter
  from the key-value store, increase it by one and then write it back
  (at the same key).  This will introduce racing condition among the
  front-ends.
- **Not handling errors**. A tribbler service call might require
  several RPC calls to the backend. It is important to properly handle
  *any* error returned by these calls.
- **Sorting by the timestamps first**. Again, Tribble Order means that
  the logic clock is the first field to consider on sorting. 
- **Use the clock argument from the front-end for the clock field of a
  new Tribble**. Well, technically, you can do that in your code
  internally as long as you can satisfy the ordering requirements
  speficied for `Home()` and `Tribs()` (you might find it very hard).
  Nonetheless, intuitively, the clock argument tells the *oldest*
  tribble a user have seen (which might be 0 if the user has not seen
  any tribble yet), hence the new posted tribble seems to better have
  a clock value that is larger than the argument.
- **Generate the clock from the timestamp**. While 64-bit can cover a
  very wide time range even in the unit of nanoseconds, you should
  keep in mind that the front-ends are running on different servers
  with arbitrary physical time differences, so it is not wise to
  generate the logical *clock* from the physical *time*.
- **Not handling old tribbles**. Note that only the most recent 100
  tribbles of a user matter. Not handling old tribbles might lead to
  worse and worse performance over time and eventually break the
  performance promise.


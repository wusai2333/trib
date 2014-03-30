## Machines

We have set up a cluster of 10 machines. You should use those for
all of the lab assignements.

- c08-11.sysnet.ucsd.edu
- c08-12.sysnet.ucsd.edu
- c08-13.sysnet.ucsd.edu
- c08-14.sysnet.ucsd.edu
- c08-15.sysnet.ucsd.edu
- c08-16.sysnet.ucsd.edu
- c08-17.sysnet.ucsd.edu
- c08-18.sysnet.ucsd.edu
- c08-19.sysnet.ucsd.edu

## Programming Language

You will write the labs in Google's [golang](http://golang.org).  It
is a young language with a language syntax at somewhere between C/C++
and Python. It comes with a very rich standard library, and also
language-level support for light-weight concurrent semantics like *go
routines* and *channels*.

Here is some key documentation on the language:

- [Go Language Documentation Page](http://golang.org/doc/)
- [Effective Go](http://golang.org/doc/effective_go.html)
- [Go Language Spec](http://golang.org/ref/spec)

While you should be able to find a lot of document about Go language
on the web, especially from the official site. If you know C++
already, here are some hints that might help you bootstrap.

- Go language code is organized in many separate *packages*.
- Different from C/C++, when defining a *variable* or *constant*, the
  *type* of it is written after the variable name.
- There are pointers in Go language, but there are no pointer
  arithmetics. For example, you cannot increase a pointer by 1, to
  point the next element in memory.
- Mapping to C/C++, in Go language, there are fixed length *arrays*.
  However, arrays are not very commonly used.  For most of the time,
  people use *slices*, which is a sliced view of an underlying array
  (often declared implicited).
- In go language, *maps* are built-in hash-based dictionaries.
- In go language, a function can have multiple return values.
- Exceptions are called `panic` and `recover`. However it is not
  encouraged to use that for error handling.
- `for` is the only loop keyword.
- *Foreach* is implemented with `range` keyword.
- Semicolons at the end of statements are optional.
- On the other hand though, trailing comma in a list is a must.
- Go language has garbage collection. It is type safe and pointer
  safe. When you have a pointer, the content it points to is always
  valid.
- Mapping to C++ concepts, in Go language, identifier that starts with
  a capital letter is *public* and visible to other modules; others
  are *private* and only visible inside modules.
- Mapping to C++ concepts, in Go language, *inheritance* is done by
  composition of anonymous members.
- Mapping to C++ concepts, in Go lanugage, virtual functions and bind
  via *interfaces*. Unlike Java, *interface* does not require explicit
  binding (via *implements* keyword). Instead, *interfaces* and bind
  on runtime. As a result, it is okay to write the implementation
  first and declare the interface afterwards.
- In Go language, you cannot have circular dependency on package
  imports.

## The Tribbler Story

Believe it or not, here is the story: some cowboy programmer wrote a
simple online microblogging service called Tribbler, and leveraging
the Web, it becomes quite popular. However, the program runs as a
single process, hence it does not scale, cannot support many
concurrent connections, and is not fault-tolerant.  Knowing that you
are taking the distributed computing system course at UCSD, he asks
you to help him. You answered his call: your goal is to refactor 
Tribbler into a distributed system, make it robust and scalable.

## Getting Started

The Tribbler project is written in golang and stored in a git
repository now. To get started, run these commands in command line:

```
$ cd                       # go to your home directory
$ mkdir -p gopath/src      # the path you use for storing golang src
$ cd gopath/src
$ git clone /classes/cse223b/sp14/labs/trib -b lab1 --depth=1
$ git clone /classes/cse223b/sp14/labs/triblab -b lab1 --depth=1
$ export GOPATH=~/gopath
$ go install ./...
```

Do some basic testing see if the framework is in good shape:

```
$ go test ./trib/...
```

Now The basic Tribbler service should be already installed on
the system in your home directory. Let's give it a try:

```
$ ~/gopath/bin/trib-front -init -addr=:rand
```

The program should show a log that it serves on a port.

Now open your browser and type in the address. For example, if the
machine you logged in is `c08-11.sysnet.ucsd.edu`, and the service is
shown running on port 27944, then open
`http://c08-11.sysnet.ucsd.edu:27944`.  You should see a list.  of
usernames, where you can view their tribs and login as them (with no
authentication). This is how the Tribbler service looks like to the
users. It is a single Web page the performs AJAX calls (a type of RPC
that is commonly used in Web 2.0) to the web server behind. The
webserver then in turn calls the Tribbler logic functions implemented
by the cowboy and returns the results back to the Web page in the
browser.

You might find it difficult to access the lab machines outside UCSD
campus. For that, you need to setup a UCSD VPN or ssh tunnel.

## Source Code Organization

The source code in the `trib` package repository is organized as follow:

- `trib` defines the high-level Tribbler logic interfaces and and
  common data structures.
- `trib/tribtest` provides several basic test cases for the
  interfaces.
- `trib/cmd/trib-front` is the web-server launcher that you just run.
- `trib/cmd/trib-back` will be the back-end storage server launcher.
- `trib/entries` defines helper functions on constructing a Tribbler
  front-end or a back-end.
- `trib/ref` is a reference implementation of the interface
  `trib.Server` interface. This is what use actually just tried via
  the Web.
- `trib/randaddr` provides helpers that generates a random port
  number.
- `trib/store` contains an in-memory thread-safe implementation of the
  Store interface. We will use this as the very basic building block
  for our back-end storage system.
- `trib/www` contains the static files (html, css, js, etc.) for the
  web front-end.

Don't be scared by the number of packages. Most of the packages are
very small. In fact, all Go language files under `trib` directory is
less than 1500 lines in total (the beauty of Go!).

Through the entire lab, you do not need to modify anything in this
`trib` repository. If you feel that you have to change some code to
complete your lab, please discuss with the TA. You are always welcome
to read the code in `trib` repository. If you find any bug and
reported it, you might get some bonus credit.

## Your Job

Your job is to complete the implementation of the `triblab` package in
the second repo that we checked out.

It would be a good practice for you to periodically commit your code
into your `triblab` git repo. Only files in that repo will be
submitted for grading.  

## Lab Roadmap

- **Lab 1**. Reimplement the Tribbler service, split the logic into a
  stateless front-end and a key-value pair back-end. The front-end
  will call the back-end via RPC. After this, we should have a
  scalable front-end design, where they can serve on multiple
  addresses concurrently. The back-end will still be single by the end
  of this lab.
- **Lab 2**. We scale up the key-value pair back-ends by consistent
  hashing the keys. The challenge for this part would be how to sort
  the Tribbles that are stored on many servers in a reasonable order
  without using a global clock.
- **Lab 3**. We make the back-end fault-tolerent, by using distributed
  hash table and replications. As a result, at the end of this lab,
  back-end servers can now join, leave, or be killed.

By the end of the lab, we will have a new Tribller service
architecture that is scalable and fault-tolerant.

## Misc

For convenience, you might set environment variables in your `.bashrc`:

```
export GOPATH=$HOME/gopath
export PATH=$PATH:$GOPATH/bin
```

We should have Vim and Emaces installed on the machines. If you need
to install other utility packages, ask the TA.

## Ready?

If you still feel comfortable with the setup,
go forward and read [Lab1](./lab1.html).


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
- c08-20.sysnet.ucsd.edu

## Programming Language

You will write the labs in Google's [golang](http://golang.org).
It is a language at somewhere between C/C++ and Python. It comes
with a very rich standard library, and also language-level support
for light-weight concurrent *go routines* and *channels*.

You should be able to find a lot of document about Go language
on the web. If you know C++ already, here are a set of hints that
might help you bootstrap.

- Code is organized in many *packages*.
- Different from C/C++, when defining a *variable* or *constant*,
  the *type* comes after the variable name.
- There are pointers in Go language, but there are no pointer
  arithmetics. For example, you cannot increase a pointer by 1,
  to point the next element in memory.
- Mapping to C/C++, in Go language, there are *arrays*, which
  are fixed length. However, arrays are not used very common.
  For most of the time, people use *slices*, which is a
  sliced view of an underlying (often implicitly declared) array.
- In go language, there are *maps*, which built-in support
  of a hash map.
- In go language, there could be multiple return values for a
  function.
- *For* is the only loop keyword.
- *Foreach* is implemented with *range* keyword.
- Semicolons at the end of statements are optional.
- On the other hand though, trailing comma in a list is a must.
- Go language has garbage collection. It is type safe and
  pointer safe. When you have a pointer, the content it points
  to is always valid.
- Mapping to C++ concepts, in Go language,
  identifier that starts with a capital letter is *public* and
  visible to other modules; others are *private* and only visible
  inside modules.
- Mapping to C++ concepts, in Go language,
  *inheritance* is done by composition of anonymous members.
- Mapping to C++ concepts, in Go lanugage,
  virtual functions and bind via *interfaces*. Unlike Java,
  *interface* does not require explicit binding
  (via *implements* keyword). Instead, *interfaces* and bind
  on runtime. As a result, it is okay to write the implementation
  first and declare the interface afterwards.
- In Go language, you cannot have circular dependency on package
  imports.

For more documentation on the language:

- [Go Language Documentation Page](http://golang.org/doc/)
- [Effective Go](http://golang.org/doc/effective_go.html)
- [Go Language Spec](http://golang.org/ref/spec)

## The Tribbler Story

The project's name is Tribbler. Here is the (fake) story: some cowboy
programmer in the wild wrote a simple online microblogging service
called Tribbler, and leveraging the world-wide-web,
it becomes quite popular. However, the program is written as a single
process program, and hence it does not scale,
cannot support many concurrent connections, and is of course
not fault-tolerant. Knowing that you are
going to take a distributed computing system course at UCSD, he
asks you to help him. Taking the challenge, your goal is to refactor
the program into a distributed system, make it robust and scalable.

## Getting Started

The Tribbler project is written in golang and stored in a git
repository now. To get start, run these commands:

```
$ cd                       # go to your home directory
$ mkdir -p gopath/src      # the path you use for storing golang src
$ export GOPATH=~/gopath
$ cd gopath/src
$ git clone /classes/xxxxxx/lab/trib -b lab1 -depth=1
$ git clone /classes/xxxxxx/lab/triblab -b lab1 -depth=1
$ go install ./...
```

The basic Tribbler service should be already installed on
the system in your home directory. Let's give it a try:

```
$ ~/gopath/bin/trib-front -init -addr=:rand
```

The program should show a log that it serves on a port.

Now open your browser and type in the address. For example,
if the machine you logged in is `c08-11.sysnet.ucsd.edu`,
and the service is shown running on port 27944,
then open `http://c08-11.sysnet.ucsd.edu:27944`.
You should see a list.
of usernames, where you can view their tribs and login as
them (with no authentication). This is how the Tribbler
service looks like to the users. It is a single Web page
the performs AJAX calls
(a type of RPC that is commonly used in Web 2.0) to the
web server behind. The webserver then in turn calls the
Tribbler logic functions implemented by the cowboy and
returns the results back to the Web page in the browser.

## Source Code Organization

The source code in the `trib` package repository is organized as follow:

- `trib` defines the high-level Tribbler logic interfaces and
  and common data structures.
- `trib/tribtest` provides several basic test cases for the
  interfaces.
- `trib/cmd/trib-front` is the web-server launcher that you just run.
- `trib/cmd/trib-back` will be the back-end storage server
launcher.
- `trib/entries` defines helper functions on constructing a
  Tribbler front-end or a back-end.
- `trib/ref` is a reference implementation of the interface
  `trib.Server`
  interface. This is what use actually just tried via the Web.
- `trib/randaddr` provides helpers that generates a random port
  number.
- `trib/store` contains an in-memory thread-safe
  implementation of the Store interface. We will use this
  as the very basic building block for our back-end storage system.
- `trib/www` contains the static files (html, css, js, etc.) for
  the web front-end.


Don't be scared by the number of packages. Most of the packages
are very small. In fact, all Go language files under `trib`
directory is less than 1500 lines in total (the beauty of Go!).

Through the entire lab, you should not need to modify anything
in this repository. If you feel that you have to change something,
please discuss with the TA. You are always welcome to read the code
in `trib` repository. If you find any bug and reported it,
you will get some bonus credit.

## Your Job

Your job is to complete the implementation of the `triblab` package
in the second repo that we checked out.

You should always write your code for the labs in the `triblab` repo,
and the is the only repo that you will submit for grading.
In fact, for most of the time, you might be just working on the machines
with `triblab` as your working directory.

For convenience, you might set environment variables in your `.bashrc`:

```
export GOPATH=$HOME/gopath
export PATH=$PATH:$GOPATH/bin
```

We should have Vim and Emaces installed on the machines. If you need
install other utility packages, ask the TA.

## Ready?

If you still feel comfortable with the setup,
go forward and read [Lab1](./lab1.html).

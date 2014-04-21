## Lab3

Welcome to Lab3. The goal of this lab is to take the bin storage that
we implemented in Lab2 and make it fault-tolerant.

## Get Your Repo Up-to-date

```
$ cd ~/gopath/src/trib
$ git pull /classes/cse223b/sp14/labs/trib lab3
$ cd ~/gopath/src/triblab
$ git pull /classes/cse223b/sp14/labs/triblab lab3
```

Not many changes, only some small things, should be painless.

## System and Failure Model

There could be up to 1000 back-ends. Back-ends may join and leave at 
any time, but there will be at least 1 back-end online. Also, you can
assume that each back-end join or leave event will have a time
interval of 30 seconds in between, and this time duration will be
enough for you to migrate storage.

There will be at least 3 keepers. Keepers may join and leave at 
any time, but there will be at least 1 keeper online. Also, you can
assume that each keeper join or leave event will have a time interval
of 1 minute in between.

## Consistency Model

To tolerate failures, you have to save the data of each key on
multiple places, and we will have a slightly relaxed consistency
model.

- `Clock()` will still be the same.

## Turning In

First, make sure that you have committed every piece of your code into
the repository `triblab`. Then just type `make turnin` under the root
of your repository. It will generate a `turnin.zip` that contains
everything in your gitt repo, and will then copy the zip file to a
place where only the lab instructors can read.

## Happy Lab3

In this directory you will find utilities related to this library

([github.com/generaltso/linguist](https://github.com/generaltso/linguist).)

## l/

A reference implementation.

## generate\_classifier/
## generate\_static/

These programs generate the data used by linguist.

`$ cd cmd`

And then use `./generate-classifier` and `./generate-static`.

### Prerequisites

 - gofmt

 - [github.com/jteeuwen/go-bindata](https://github.com/jteeuwen/go-bindata)

## tokenizer_test/

A small program to test the tokenizer.

Usage:
```
$ cd $GOPATH/src/github.com/generaltso/linguist/cmd/tokenizer_test
$ go run main.go /path/to/some/file
```

Note you can't run `go run main.go main.go` <s>for some reason</s> because the go tool will greedily accept any .go files as source files rather than arguments, <s>unless you</s> so compile it first with `go build` to run on .go files.

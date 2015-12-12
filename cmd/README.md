In this directory you will find utilities related to this library

([github.com/generaltso/linguist](https://github.com/generaltso/linguist).)

## l/

A reference implementation.

## generate\_classifier/
## generate\_static/

These programs generate the data used by linguist.

Use `./generate-classifier` and `./generate-static` rather than running them directly.

## tokenizer_test/

A small program to test the tokenizer.

Usage:
```
$ cd $GOPATH/src/github.com/generaltso/linguist/cmd/tokenizer_test
$ go run main.go /path/to/some/file
```

Note you can't run go run main.go on main.go itself for some reason, unless you compile it first with `go build`

/*
Detect programming language of source files.
Go port of GitHub Linguist: https://github.com/github/linguist

Prerequisites:

    go get github.com/jteeuwen/go-bindata/go-bindata

Installation:

    go get -d github.com/generaltso/linguist
    cd $GOPATH/src/github.com/generaltso/linguist
    make

Optional:
    go install github.com/generaltso/linguist/cmd/l


Usage:

Please refer to the source code for the reference implementation at:

https://github.com/generaltso/linguist/tree/master/cmd/l


See also:

https://github.com/generaltso/linguist/tree/master/tokenizer
*/
package linguist

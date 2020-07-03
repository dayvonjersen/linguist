/*
Detect programming language of source files.
Go port of GitHub Linguist: https://github.com/github/linguist

Prerequisites:

    go get github.com/jteeuwen/go-bindata/go-bindata

Installation:

    mkdir -p $GOPATH/src/github.com/dayvonjersen/linguist
    git clone --depth=1 https://github.com/dayvonjersen/linguist $GOPATH/src/github.com/dayvonjersen/linguist
    go get -d github.com/dayvonjersen/linguist
    cd $GOPATH/src/github.com/dayvonjersen/linguist
    make
    l

Usage:

Please refer to the source code for the reference implementation at:

https://github.com/dayvonjersen/linguist/tree/master/cmd/l


See also:

https://github.com/dayvonjersen/linguist/tree/master/tokenizer
*/
package linguist

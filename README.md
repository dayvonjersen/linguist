# linguist

[![godoc reference](https://godoc.org/github.com/generaltso/linguist?status.png)](https://godoc.org/github.com/generaltso/linguist)

Go port of [github linguist](https://github.com/github/linguist).

Many thanks to [@petermattis](https://github.com/petermattis) for his initial work in laying the groundwork of creating this project, and especially for suggesting the use of naive Bayesian classification.

Thanks also to [@jbrukh](https://github.com/jbrukh) for [github.com/jbrukh/bayesian](https://github.com/jbrukh/bayesian)

# install

### prerequisites:

```
go get github.com/jteeuwen/go-bindata/go-bindata
```

```
mkdir -p $GOPATH/src/github.com/generaltso/linguist
git clone --depth=1 https://github.com/generaltso/linguist $GOPATH/src/github.com/generaltso/linguist
go get -d github.com/generaltso/linguist
cd $GOPATH/src/github.com/generaltso/linguist
make
l
```

## see also

[command-line reference implentation](cmd/l) which is documented separately

[tokenizer](tokenizer/tokenizer.go) | ([godoc reference](https://godoc.org/github.com/generaltso/linguist/tokenizer))

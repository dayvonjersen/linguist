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
go get -d github.com/generaltso/linguist
cd $GOPATH/src/github.com/generaltso/linguist
make
```

### optional

```
go install github.com/generaltso/linguist/cmd/l

```

[the command-line reference implentation](cmd/l) is documented separately

## see also

[the tokenizer I made for this project](tokenizer/tokenizer.go) which you can import:

```
import "github.com/generaltso/linguist/tokenizer"
```

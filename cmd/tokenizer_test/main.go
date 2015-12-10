package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/generaltso/linguist/tokenizer"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	for _, file := range flag.Args() {
		fmt.Println("attempting to read", file)
		fp, e := filepath.Abs(file)
		checkErr(e)
		f, e := os.Open(fp)
		checkErr(e)
		meta, e := ioutil.ReadAll(f)
		checkErr(e)
		fmt.Printf("%#v\n", tokenizer.Tokenize(meta))
	}

}

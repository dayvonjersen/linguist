// +build ignore

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/generaltso/linguist/tokenizer"
	"github.com/jbrukh/bayesian"
)

var training_dir string

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	flag.StringVar(&training_dir, "training-dir", "/tmp/linguist/samples", `path to a directory of programming language samples
from which to train the Naive Bayesian Classifier

this program defaults to using the samples distributed with
https://github.com/github/linguist

you can use the default setting for this flag and thereby omit it by:

    git clone git@github.com:github/linguist /tmp/linguist
`)

	flag.Parse()

	log.SetOutput(os.Stderr)

	classes := make([]bayesian.Class, 1)
	documents := make(map[bayesian.Class][]string)

	log.Println("Scanning", training_dir, "...")

	tdir, err := os.Open(training_dir)
	checkErr(err)

	ldirs, err := tdir.Readdir(-1)
	checkErr(err)

	for _, ldir := range ldirs {
		lang := ldir.Name()
		if !ldir.IsDir() {
			log.Println("unexpected file:", lang)
			continue
		}

		log.Println("Found Language:", lang)

		var class bayesian.Class = bayesian.Class(lang)
		classes = append(classes, class)

		w := []string{}

		sample_dir := training_dir + "/" + lang

		sdir, _ := os.Open(sample_dir)
		files, err := sdir.Readdir(-1)
		checkErr(err)
		for _, file := range files {
			fp := sample_dir + "/" + file.Name()
			if file.IsDir() {
				log.Println("Skipping subdirectory", fp, "...")
				continue
			}
			f, err := os.Open(fp)
			checkErr(err)
			contents, err := ioutil.ReadAll(f)
			checkErr(err)
			w = append(w, tokenizer.Tokenize(contents)...)
		}
		documents[class] = w
	}

	log.Println("Creating bayesian.Classifier ...")
	clsf := bayesian.NewClassifier(classes...)
	for cls, dox := range documents {
		clsf.Learn(dox, cls)
	}

	log.Println("Serializing and exporting bayesian.Classifier into data/classifier ...")
	checkErr(clsf.WriteToFile("data/classifier"))

	log.Println("Done.")
}

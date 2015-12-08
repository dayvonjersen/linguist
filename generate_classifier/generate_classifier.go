package main

import (
	"encoding/json"
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
	log.Println("Opening ../samples.json ...")
	f, err := os.Open("../samples.json")
	checkErr(err)
	log.Println("Reading ../samples.json ...")
	b, err := ioutil.ReadAll(f)
	checkErr(err)

	log.Println("Unmarshaling ../samples.json ...")
	var dat map[string]interface{}
	checkErr(json.Unmarshal(b, &dat))

	tokens := dat["tokens"].(map[string]interface{})

	classes := make([]bayesian.Class, 1)
	documents := make(map[bayesian.Class][]string)

	log.Println("Reading data into vars ...")
	for key, _ := range tokens {
		if key == "" {
			continue
		}
		//log.Println("Found Language:", key)

		var class bayesian.Class = bayesian.Class(key)

		classes = append(classes, class)

		w := []string{}
		sample_dir := training_dir + "/" + key
		dir, err := os.Open(sample_dir)
		if os.IsNotExist(err) {
			log.Println("DIRECTORY NOT FOUND:", sample_dir)
		} else {
			files, err := dir.Readdir(-1)
			checkErr(err)
			for _, file := range files {
				if file.IsDir() {
					log.Println("Skipping subdirectory", sample_dir+"/"+file.Name(), "...")
					continue
				}
				fp := sample_dir + "/" + file.Name()
				//log.Println("Tokenizing", fp, "...")
				f, err := os.Open(fp)
				checkErr(err)
				contents, err := ioutil.ReadAll(f)
				checkErr(err)
				w = append(w, tokenizer.Tokenize(contents)...)
			}
		}
		documents[class] = w
	}

	log.Println("Creating bayesian.Classifier ...")
	clsf := bayesian.NewClassifier(classes...)
	for cls, dox := range documents {
		clsf.Learn(dox, cls)
	}

	log.Println("Serializing and exporting bayesian.Classifier into ./out/classifier ...")
	checkErr(clsf.WriteToFile("./out/classifier"))

	log.Println("Done.")
}

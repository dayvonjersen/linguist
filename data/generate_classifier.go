// +build ignore

/*
   This program trains a naive bayesian classifier
   provided by https://github.com/jbrukh/bayesian
   on a set of source code files
   provided by https://github.com/github/linguist

   This file is meant by run by go generate,
   refer to generate.go for its intended invokation
*/
package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/generaltso/linguist/tokenizer"
	"github.com/jbrukh/bayesian"
)

func main() {
	const (
		sourcePath = "./linguist/samples"
		outfile    = "./classifier"
		quiet      = false
	)

	if quiet {
		log.SetOutput(ioutil.Discard)
	}

	// the sample files are stored into structs
	type sampleFile struct {
		lang, fp string
		tokens   []string
	}

	// first we only read all the paths of the sample files
	// and their corresponding and language names into:
	sampleFiles := []*sampleFile{}

	// and store all the language names into:
	languages := []string{}

	/*
			   github/linguist has directory structure:

			   ...
			   ├── samples
			   │   ├── (name of programming language)
			   │   │   ├── (sample file in language)
			   │   │   ├── (sample file in language)
			   │   │   └── (sample file in language)
			   │   ├── (name of another programming language)
			   │   │   └── (sample file)
			   ...

		       the following hard-coded logic expects this layout
	*/

	log.Println("Scanning", sourcePath, "...")
	srcDir, err := os.Open(sourcePath)
	checkErr(err)

	subDirs, err := srcDir.Readdir(-1)
	checkErr(err)

	for _, langDir := range subDirs {
		lang := langDir.Name()
		if !langDir.IsDir() {
			log.Println("unexpected file:", lang)
			continue
		}
		log.Println("Found Language:", lang)
		languages = append(languages, lang)

		samplePath := sourcePath + "/" + lang
		sampleDir, _ := os.Open(samplePath)
		files, err := sampleDir.Readdir(-1)
		checkErr(err)
		for _, file := range files {
			fp := samplePath + "/" + file.Name()
			if file.IsDir() {
				log.Println("Skipping subdirectory", fp, "...")
				continue
			}
			sampleFiles = append(sampleFiles, &sampleFile{lang, fp, nil})
		}
	}

	log.Println("Parsing files...")
	// then we concurrently read and tokenize the samples
	sampleChan := make(chan *sampleFile)
	readyChan := make(chan struct{})
	received := 0
	// and store them in a map of languages as keys and slices of tokens as values
	dox := map[string][]string{}
	// (receives the processed files and stores their tokens with their language)
	go func() {
		for {
			s := <-sampleChan
			if doc, ok := dox[s.lang]; ok {
				dox[s.lang] = append(doc, s.tokens...)
			} else {
				dox[s.lang] = s.tokens
			}
			received++
			if received == len(sampleFiles) {
				close(readyChan)
				return
			}
		}
	}()

	// (concurrently reads and tokenizes files)
	for _, s := range sampleFiles {
		go func() {
			f, err := os.Open(s.fp)
			checkErr(err)
			contents, err := ioutil.ReadAll(f)
			f.Close()
			checkErr(err)
			s.tokens = tokenizer.Tokenize(contents)
			sampleChan <- s
		}()
	}

	// once that's done
	<-readyChan
	// we train the classifier in the arbitrary manner that its API demands
	classes := make([]bayesian.Class, 1)
	documents := make(map[bayesian.Class][]string)
	for _, lang := range languages {
		var class = bayesian.Class(lang)
		classes = append(classes, class)
		documents[class] = dox[lang]
	}
	log.Println("Creating bayesian.Classifier ...")
	clsf := bayesian.NewClassifier(classes...)
	for cls, dox := range documents {
		clsf.Learn(dox, cls)
	}

	// and write the data to disk
	log.Println("Serializing and exporting bayesian.Classifier to", outfile, "...")
	checkErr(clsf.WriteToFile("classifier"))

	log.Println("Done.")
}
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

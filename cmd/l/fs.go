package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/generaltso/linguist"
	"github.com/lintianzhi/ignore"
)

var isIgnored func(string) bool

func initGitIgnore() {
	if fileExists(".git") && fileExists(".gitignore") {
		log.Println("found .git directory and .gitignore")
		gitIgn, err := ignore.NewGitIgn(".gitignore")
		checkErr(err)
		gitIgn.Start(".")
		isIgnored = func(filename string) bool {
			return gitIgn.TestIgnore(filename)
		}
	} else {
		log.Println("no .gitignore found")
		isIgnored = func(filename string) bool {
			return false
		}
	}
}

// shoutouts to php
func fileGetContents(filename string) []byte {
	log.Println("reading contents of", filename)
	contents, err := ioutil.ReadFile(filename)
	checkErr(err)
	return contents
}
func fileExists(filename string) bool {
	log.Println("opening file", filename)
	g, err := os.Open(filename)
	g.Close()
	if os.IsNotExist(err) {
		log.Println(filename, "does not exist")
		return false
	}
	checkErr(err)
	return true
}

func processDir(dirname string) {
	cwd, err := os.Open(dirname)
	checkErr(err)
	files, err := cwd.Readdir(0)
	checkErr(err)
	checkErr(os.Chdir(dirname))
	for _, file := range files {
		name := file.Name()
		size := int(file.Size())
		log.Println("with file: ", name)
		if size == 0 {
			log.Println(name, "is empty file, skipping")
			continue
		}
		if isIgnored(dirname + string(os.PathSeparator) + name) {
			log.Println(name, "is ignored, skipping")
			continue
		}
		if file.IsDir() {
			if name == ".git" {
				log.Println(".git directory, skipping")
				continue
			}
			processDir(name)
		} else {
			if linguist.IsVendored(name) {
				log.Println(name, "is vendored, skipping")
				continue
			}

			by_name := linguist.DetectFromFilename(name)
			if by_name != "" {
				log.Println(name, "got result by name: ", by_name)
				putResult(by_name, size)
				continue
			}

			contents := fileGetContents(name)

			if linguist.IsBinary(contents) {
				log.Println(name, "is (likely) binary file, skipping")
				continue
			}

			by_data := linguist.DetectFromContents(contents)
			if by_data != "" {
				log.Println(name, "got result by data: ", by_data)
				putResult(by_data, size)
				continue
			}

			log.Println(name, "got no result!!")
			putResult("(unknown)", size)
		}
	}
	checkErr(os.Chdir(".."))
}

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
		log.Println("with file: ", file.Name())
		//abs, err := filepath.Abs(file.Name())
		//checkErr(err)
		size := int(file.Size())
		if size == 0 {
			log.Println("size is 0 for", file.Name())
			continue
		}
		if isIgnored(dirname + string(os.PathSeparator) + file.Name()) {
			log.Println(file.Name(), "is ignored")
			continue
		}
		if file.IsDir() {
			if file.Name() == ".git" {
				log.Println("skipping .git directory")
				continue
			}
			processDir(file.Name())
		} else if !linguist.IsVendored(file.Name()) {
			by_name, shouldIgnore := getLangFromFilename(file.Name())
			if shouldIgnore {
				log.Println("DetectMimeFromFilename says to ignore type: ", by_name)
				log.Println("Ignoring", file.Name())
				continue
			} else if by_name != "" {
				log.Println("got result by name: ", by_name)
				putResult(by_name, size)
				continue
			}

			by_data, shouldIgnore := getLangFromContents(fileGetContents(file.Name()))
			if shouldIgnore {
				log.Println("DetectMimeFromContents says to ignore type: ", by_data)
				log.Println("Ignoring", file.Name())
				continue
			} else if by_data != "" {
				log.Println("got result by data: ", by_data)
				putResult(by_data, size)
				continue
			}

			log.Println("got no result for: ", file.Name())
			putResult("(unknown)", size)
		} else {
			log.Println(file.Name(), "is vendored")
		}
	}
	checkErr(os.Chdir(".."))
}

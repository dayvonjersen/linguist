package main

import (
	"io/ioutil"
	"os"

	"github.com/generaltso/linguist"
	"github.com/lintianzhi/ignore"
)

var isIgnored func(string) bool

func initGitIgnore() {
	if fileExists(".gitignore") {
		gitIgn, err := ignore.NewGitIgn(".gitignore")
		checkErr(err)
		gitIgn.Start(".")
		isIgnored = func(filename string) bool {
			return gitIgn.TestIgnore(filename)
		}
	} else {
		isIgnored = func(filename string) bool {
			return false
		}
	}
}

// shoutouts to php
func fileGetContents(filename string) []byte {
	contents, err := ioutil.ReadFile(filename)
	checkErr(err)
	return contents
}
func fileExists(filename string) bool {
	g, err := os.Open(filename)
	g.Close()
	if os.IsNotExist(err) {
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
		//abs, err := filepath.Abs(file.Name())
		//checkErr(err)
		size := int(file.Size())
		if size == 0 {
			continue
		}
		if isIgnored(dirname + string(os.PathSeparator) + file.Name()) {
			continue
		}
		if file.IsDir() {
			if file.Name() == ".git" {
				continue
			}
			processDir(file.Name())
		} else if !linguist.IsVendored(file.Name()) {
			by_name := getLangFromFilename(file.Name())
			if by_name != "" {
				putResult(by_name, size)
				continue
			}

			by_data := getLangFromContents(fileGetContents(file.Name()))
			if by_data != "" {
				putResult(by_data, size)
				continue
			}

			putResult("(unknown)", size)
		}
	}
	checkErr(os.Chdir(".."))
}

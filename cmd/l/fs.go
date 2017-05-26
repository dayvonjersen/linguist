package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/generaltso/linguist"
)

var isIgnored func(string) bool

func initGitIgnore() {
	if fileExists(".git") && fileExists(".gitignore") {
		log.Println("found .git directory and .gitignore")

		f, err := os.Open(".gitignore")
		checkErr(err)

		pathlist, err := ioutil.ReadAll(f)
		checkErr(err)

		ignore := []string{}
		except := []string{}
		for _, path := range strings.Split(string(pathlist), "\n") {
			path = strings.TrimSpace(path)
			if len(path) == 0 || string(path[0]) == "#" {
				continue
			}
			isExcept := false
			if string(path[0]) == "!" {
				isExcept = true
				path = path[1:]
			}
			fields := strings.Split(path, " ")
			p := fields[len(fields)-1:][0]
			p = strings.Trim(p, string(filepath.Separator))
			if isExcept {
				except = append(except, p)
			} else {
				ignore = append(ignore, p)
			}
		}
		isIgnored = func(filename string) bool {
			for _, p := range ignore {
				if m, _ := filepath.Match(p, filename); m {
					for _, e := range except {
						if m, _ := filepath.Match(e, filename); m {
							return false
						}
					}
					return true
				}
			}
			return false
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

	// read only first 512 bytes of files
	contents := make([]byte, 512)
	f, err := os.Open(filename)
	checkErr(err)
	_, err = f.Read(contents)
	f.Close()
	if err != io.EOF {
		checkErr(err)
	}
	return contents
}

func fileExists(filename string) bool {
	log.Println("opening file", filename)
	f, err := os.Open(filename)
	f.Close()
	if os.IsNotExist(err) {
		log.Println(filename, "does not exist")
		return false
	}
	checkErr(err)
	return true
}

func processDir(dirname string) {
	filepath.Walk(dirname, func(path string, file os.FileInfo, err error) error {
		size := int(file.Size())
		log.Println("with file: ", path)
		log.Println(path, "is", size, "bytes")
		if size == 0 {
			log.Println(path, "is empty file, skipping")
			return nil
		}
		if !unignore_filenames && isIgnored(path) {
			log.Println(path, "is ignored, skipping")
			ignored_paths++
			if file.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if file.IsDir() {
			if file.Name() == ".git" {
				log.Println(".git directory, skipping")
				return filepath.SkipDir
			}
		} else if (file.Mode() & os.ModeSymlink) == 0 {
			if !unignore_filenames && linguist.ShouldIgnoreFilename(path) {
				log.Println(path, ": filename should be ignored, skipping")
				ignored_paths++
				return nil
			}

			by_name := linguist.LanguageByFilename(path)
			if by_name != "" {
				log.Println(path, "got result by name: ", by_name)
				putResult(by_name, size)
				return nil
			}

			contents := fileGetContents(path)

			if !unignore_contents && linguist.ShouldIgnoreContents(contents) {
				log.Println(path, ": contents should be ignored, skipping")
				ignored_paths++
				return nil
			}

			hints := linguist.LanguageHints(path)
			log.Printf("%s got language hints: %#v\n", path, hints)
			by_data := linguist.LanguageByContents(contents, hints)

			if by_data != "" {
				log.Println(path, "got result by data: ", by_data)
				putResult(by_data, size)
				return nil
			}

			log.Println(path, "got no result!!")
			putResult("(unknown)", size)
		}
		return nil
	})
}

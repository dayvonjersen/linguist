package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/generaltso/linguist"
)

var isIgnored func(string) bool
var isDetectedInGitAttributes func(string) string

func initLinguistAttributes() {
	ignore := []string{}
	except := []string{}
	detected := make(map[string]string)

	if fileExists(".gitignore") {
		log.Println("found .gitignore")

		f, err := os.Open(".gitignore")
		checkErr(err)
		defer f.Close()

		ignoreScanner := bufio.NewScanner(f)
		for ignoreScanner.Scan() {
			var isExcept bool
			path := strings.TrimSpace(ignoreScanner.Text())
			// if it's whitespace or a comment
			if len(path) == 0 || string(path[0]) == "#" {
				continue
			}
			if string(path[0]) == "!" {
				isExcept = true
				path = path[1:]
			}
			p := strings.Trim(path, string(filepath.Separator))
			if isExcept {
				except = append(except, p)
			} else {
				ignore = append(ignore, p)
			}
		}
		checkErr(ignoreScanner.Err())
	}

	if fileExists(".gitattributes") {
		log.Println("found .gitattributes")

		f, err := os.Open(".gitattributes")
		checkErr(err)
		defer f.Close()

		attributeScanner := bufio.NewScanner(f)
		var lineNumber int
		for attributeScanner.Scan() {
			lineNumber++
			line := strings.TrimSpace(attributeScanner.Text())
			words := strings.Fields(line)
			if len(words) != 2 {
				log.Printf("invalid line in .gitattributes at L%d: '%s'\n", lineNumber, line)
				continue
			}
			path := strings.Trim(words[0], string(filepath.Separator))
			attribute := words[1]
			if strings.HasPrefix(attribute, "linguist-documentation") || strings.HasPrefix(attribute, "linguist-vendored") || strings.HasPrefix(attribute, "linguist-generated") {
				if strings.HasSuffix(strings.ToLower(attribute), "false") {
					except = append(except, path)
				}
			} else if strings.HasPrefix(attribute, "linguist-language") {
				attr := strings.Split(attribute, "=")
				if len(attr) != 2 {
					log.Printf("invalid line in .gitattributes at L%d: '%s'\n", lineNumber, line)
					continue
				}
				language := attr[1]
				detected[path] = language
			}
		}
		checkErr(attributeScanner.Err())
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
	isDetectedInGitAttributes = func(filename string) string {
		for p, lang := range detected {
			if m, _ := filepath.Match(p, filename); m {
				return lang
			}
		}
		return ""
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

			byGitAttr := isDetectedInGitAttributes(path)
			if byGitAttr != "" {
				log.Println(path, "got result by .gitattributes: ", byGitAttr)
				log.Println("")
				putResult(byGitAttr, size)
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

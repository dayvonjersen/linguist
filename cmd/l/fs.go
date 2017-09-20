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

//
// .gitattributes is used for overriding language detection behavior for specified filepaths
//
// the following format applies:
//
// path/to/something linguist-documentation
// # this is a comment
// path/to/something linguist-documentation=false
// path/to/something linguist-language=Go
//
// see https://github.com/github/linguist#overrides for more information
//

type override struct {
	isVendored, isDocumentation, isGenerated bool
	language                                 string
}

var getOverride func(filename string) (*override, bool)

func initGitAttributes() {
	if fileExists(".git") && fileExists(".gitattributes") {
		log.Println("found .git directory and .gitattributes")

		f, err := os.Open(".gitattributes")
		checkErr(err)
		defer f.Close()

		attributes, err := ioutil.ReadAll(f)
		checkErr(err)

		overridePaths := []string{}
		overrides := map[string]*override{}

		for _, ln := range strings.Split(string(attributes), "\n") {
			fields := strings.Fields(ln)
			if len(fields) == 2 {
				path, attr, value := fields[0], fields[1], ""

				if path[0] == '#' {
					continue
				}

				if strings.Contains(attr, "=") {
					tmp := strings.Split(fields[1], "=")
					attr, value = tmp[0], tmp[1]
				}

				o := &override{}
				switch attr {
				case "linguist-documentation":
					if value != "false" {
						o.isDocumentation = true
					}
				case "linguist-vendored":
					if value != "false" {
						o.isVendored = true
					}
				case "linguist-generated":
					if value != "false" {
						o.isGenerated = true
					}
				case "linguist-language":
					if value != "" {
						o.language = value
					}
				}
				overridePaths = append(overridePaths, path)
				overrides[path] = o
			}
		}
		// reverse list of paths as a quick and dirty hack for overriding overrides e.g.:
		//
		// docs/* linguist-documentation
		// docs/formatter.rb linguist-documentation=false
		//
		{
			last := len(overridePaths) - 1
			for i := 0; i < len(overridePaths)/2; i++ {
				overridePaths[i], overridePaths[last-i] = overridePaths[last-i], overridePaths[i]
			}
		}
		getOverride = func(filename string) (*override, bool) {
			for _, p := range overridePaths {
				if m, _ := filepath.Match(p, filename); m {
					return overrides[p], true
				}
			}
			return nil, false
		}
	} else {
		log.Println("no .gitattributes found")
		getOverride = func(filename string) (*override, bool) {
			return nil, false
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
			if o, ok := getOverride(path); ok {
				if !unignore_filenames {
					if o.isVendored || o.isDocumentation || o.isGenerated {
						log.Println(path, ": filename should be ignored, skipping (GITATTRIBUTES OVERRIDE)")
						ignored_paths++
						return nil
					}
				}
				if o.language != "" {
					by_name := o.language
					log.Println(path, "got result by name: ", by_name, "(GITATTRIBUTES OVERRIDE)")
					putResult(by_name, size)
					return nil
				}
			}
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

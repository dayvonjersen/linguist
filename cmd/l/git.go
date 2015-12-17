package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/generaltso/linguist"
)

func catfile(hash string) []byte {
	log.Println("git cat-file blob", hash)
	git := exec.Command("sh", "-c", "git cat-file blob "+hash)
	stdout, err := git.StdoutPipe()
	checkErr(err)
	c := make(chan struct{})
	r := make(chan struct{})
	blob := make([]byte, 512)
	go func() {
		git.Start()
		r <- struct{}{}
		git.Wait()
		c <- struct{}{}
		log.Println("EXITED: git cat-file blob", hash)
	}()
	go func() {
		<-r
		n, err := stdout.Read(blob)
		log.Printf("Read %d bytes from %s", n, hash)
		if err != io.EOF {
			checkErr(err)
		} else {
			log.Println("Reached EOF for", hash)
		}
		git.Process.Kill()
		c <- struct{}{}
		log.Println("KILLED: git cat-file blob", hash)
	}()
	<-c
	return blob
}

func gitcmd(args string) []byte {
	log.Println("git", args)
	git := exec.Command("sh", "-c", "git "+args)
	out, err := git.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
	}
	checkErr(err)
	return out
}

func gitcmdString(args string) string {
	return string(gitcmd(args))
}

func processTree(tree_id string, parent []string) {
	ls_tree := gitcmdString("ls-tree " + tree_id)
	for _, ln := range strings.Split(ls_tree, "\n") {
		fields := strings.Split(ln, " ")
		if len(fields) != 3 {
			continue
		}
		//fmode := fields[0]
		ftype := fields[1]
		fields = strings.Split(fields[2], "\t")
		if len(fields) != 2 {
			continue
		}
		fhash := fields[0]
		fname := fields[1]

		switch ftype {
		case "tree":
			log.Println("entering subtree", fname)
			processTree(fhash, append(parent, fname))
		case "blob":
			fname = filepath.Join(append(parent, fname)...)
			cat_size := gitcmdString("cat-file -s " + fhash)
			size, err := strconv.Atoi(strings.TrimSpace(cat_size))
			checkErr(err)

			log.Println(fname, "is", size, "bytes")
			if size == 0 {
				log.Println(fname, "is empty file, skipping")
				continue
			}

			if linguist.ShouldIgnoreFilename(fname) {
				log.Println(fname, ": filename should be ignored, skipping")
				continue
			}

			by_name := linguist.LanguageByFilename(fname)
			if by_name != "" {
				log.Println(fname, "got result by name: ", by_name)
				putResult(by_name, size)
				continue
			}

			contents := catfile(fhash)

			if linguist.ShouldIgnoreContents(contents) {
				log.Println(fname, ": contents should be ignored, skipping")
				continue
			}

			hints := linguist.LanguageHints(fname)
			log.Printf("%s got language hints: %#v\n", fname, hints)
			by_data := linguist.LanguageByContents(contents, hints)

			if by_data != "" {
				log.Println(fname, "got result by data: ", by_data)
				putResult(by_data, size)
				continue
			}
			log.Println(fname, "got no result!!")
			putResult("(unknown)", size)
		case "commit":
			log.Println(fname, "is a git submodule (ftype == \"commit\"), skipping")
			continue
		default:
			println("currently unsupported ftype:" + ftype)
		}
	}
}

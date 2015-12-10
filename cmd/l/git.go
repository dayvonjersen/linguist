package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/generaltso/linguist"
)

func gitcmd(args string) []byte {
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

func processTree(tree_id string) {
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
			processTree(fhash)
		case "blob":
			cat_size := gitcmdString("cat-file -s " + fhash)
			size, err := strconv.Atoi(strings.TrimSpace(cat_size))
			checkErr(err)

			if size == 0 {
				log.Println("omitting empty file", fname)
				continue
			}
			if linguist.IsVendored(fname) {
				log.Println("omitting vendored file", fname)
				continue
			}

			by_name, shouldIgnore := getLangFromFilename(fname)
			if shouldIgnore {
				log.Println("DetectMimeFromFilename says to ignore type: ", by_name)
				log.Println("Ignoring", fname)
				continue
			} else if by_name != "" {
				log.Println("got result by name: ", by_name)
				putResult(by_name, size)
				continue
			}

			cat_data := gitcmd("cat-file blob " + fhash)
			by_data, shouldIgnore := getLangFromContents(cat_data)
			if shouldIgnore {
				log.Println("DetectMimeFromContents says to ignore type: ", by_data)
				log.Println("Ignoring", fname)
				continue
			} else if by_data != "" {
				log.Println("got result by data: ", by_data)
				putResult(by_data, size)
				continue
			}
			log.Println("got no result for: ", fname)
			putResult("(unknown)", size)
		case "commit":
			log.Println("omitting commit", fname)
			continue
		default:
			println("currently unsupported ftype:" + ftype)
		}
	}
}

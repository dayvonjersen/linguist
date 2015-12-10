package main

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/generaltso/linguist"
)

func gitcmd(args string) []byte {
	git := exec.Command("sh", "-c", "git "+args)
	out, err := git.CombinedOutput()
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
			processTree(fhash)
		case "blob":
			cat_size := gitcmdString("cat-file -s " + fhash)
			size, err := strconv.Atoi(strings.TrimSpace(cat_size))
			checkErr(err)

			if size == 0 {
				continue
			}
			if linguist.IsVendored(fname) {
				continue
			}

			by_name := getLangFromFilename(fname)
			if by_name != "" {
				putResult(by_name, size)
				continue
			}

			cat_data := gitcmd("cat-file blob " + fhash)
			by_data := getLangFromContents(cat_data)
			if by_data != "" {
				putResult(by_data, size)
				continue
			}

			putResult("(unknown)", size)
		default:
			println("unsupported ftype:" + ftype)
		}
	}
}

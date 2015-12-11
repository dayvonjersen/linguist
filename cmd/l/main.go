package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/generaltso/linguist"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var input_mode_git bool
var input_mode_fs bool
var input_git_tree string
var output_json bool
var output_json_with_colors bool
var output_limit int
var output_debug bool

type language struct {
	Language string
	Percent  float64
}

type language_color struct {
	Language string
	Percent  float64
	Color    string
}

var langs map[string]int = make(map[string]int)
var total_size int = 0
var res map[string]int = make(map[string]int)
var num_files int = 0
var max_len int = 0

func putResult(language string, size int) {
	res[language]++
	langs[language] += size
	total_size += size
	num_files++
	if len(language) > max_len {
		max_len = len(language)
	}
}

func main() {
	flag.BoolVar(&output_debug, "debug", false, "Print debug information.")
	flag.BoolVar(&input_mode_git, "git", false, "Scan for files using git ls-tree and cat-file, rather than filesystem.")
	flag.BoolVar(&input_mode_fs, "fs", false, "Scan for files using filesystem.")
	flag.StringVar(&input_git_tree, "git-tree", "HEAD", "tree-ish root to scan. See also man git(1).")
	flag.BoolVar(&output_json, "json", false, "Output results in JSON format.")
	flag.BoolVar(&output_json_with_colors, "json-with-colors", false, "Output results in JSON format, including any HTML color codes defined for associated languages.")
	flag.IntVar(&output_limit, "limit", 10, "Limit result set to n results. n <= 0 indicates unlimited result set.")
	flag.Parse()

	output_json = output_json || output_json_with_colors

	if !output_debug {
		log.SetOutput(ioutil.Discard)
	}

	var default_input_mode_git bool
	var default_input_mode_fs bool
	if fileExists(".git") {
		default_input_mode_git = true
		default_input_mode_fs = false
	} else {
		default_input_mode_git = false
		default_input_mode_fs = true
	}

	if !input_mode_git && !input_mode_fs {
		input_mode_git = default_input_mode_git
		input_mode_fs = default_input_mode_fs
	}

	if !input_mode_git && input_git_tree != "HEAD" {
		input_mode_git = true
		input_mode_fs = false
	}

	if input_mode_git && input_mode_fs {
		fmt.Println("Please choose one of -git or -fs as flags, but not both.")
		fmt.Println("You can omit the flags to get the default behavior,")
		fmt.Printf("which for the current directory is %s\n", func() string {
			switch {
			case default_input_mode_git:
				return "git"
			case default_input_mode_fs:
				return "fs"
			}
			return "undefined"
		}())
		os.Exit(1)
	}

	if input_mode_fs {
		initGitIgnore()
		processDir(".")
	}

	if input_mode_git {
		processTree(input_git_tree)
	}

	results := []float64{}
	qqq := map[float64]string{}
	for lang, num := range langs {
		res := (float64(num) / float64(total_size)) * 100.0
		results = append(results, res)
		qqq[res] = lang
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(results)))

	if output_json {
		out := []interface{}{}
		for i, percent := range results {
			if output_limit > 0 && i >= output_limit {
				break
			}
			var l interface{}
			if output_json_with_colors {
				l = language_color{qqq[percent], percent, linguist.GetColor(qqq[percent])}
			} else {
				l = language{qqq[percent], percent}
			}
			out = append(out, l)
		}
		j, err := json.MarshalIndent(out, "", "  ")
		checkErr(err)
		fmt.Println(string(j))
		os.Exit(0)
	}
	fmtstr := fmt.Sprintf("%% %ds", max_len)
	fmtstr += ": %07.4f%%\n"

	for i, percent := range results {
		if output_limit > 0 && i >= output_limit {
			break
		}
		fmt.Printf(fmtstr, qqq[percent], percent)
	}
	fmt.Printf("\n%d languages detected in %d files\n", len(langs), num_files)
}

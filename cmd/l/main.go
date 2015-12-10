package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
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

func main() {
	var default_input_mode_git bool
	var default_input_mode_fs bool
	if fileExists(".git") {
		default_input_mode_git = true
		default_input_mode_fs = false
	} else {
		default_input_mode_git = false
		default_input_mode_fs = true
	}

	flag.BoolVar(&input_mode_git, "git", default_input_mode_git, "Scan for files using git ls-tree and cat-file, rather than filesystem.")
	flag.BoolVar(&input_mode_fs, "fs", default_input_mode_fs, "Scan for files using filesystem.")
	flag.StringVar(&input_git_tree, "git-tree", "HEAD", "tree-ish root to scan. See also man git(1).")
	flag.BoolVar(&output_json, "json", false, "Output results in JSON format.")
	flag.BoolVar(&output_json_with_colors, "json-with-colors", false, "Output results in JSON format, including any HTML color codes defined for associated languages.")
	flag.IntVar(&output_limit, "limit", 10, "Limit result set to n results. n <= 0 indicates unlimited result set.")
	flag.Parse()

	output_json = output_json || output_json_with_colors

	if input_mode_git && input_mode_fs || (input_mode_git == false && input_mode_fs == false) {
		input_mode_git = !default_input_mode_git
		input_mode_fs = !default_input_mode_fs
	}

	switch true {
	case input_mode_git:
		processTree(input_git_tree)
	case input_mode_fs:
		initGitIgnore()
		processDir(".")
	}

	fmtstr := fmt.Sprintf("%% %ds", max_len)
	fmtstr += ": %07.4f%%\n"
	results := []float64{}
	qqq := map[float64]string{}
	for lang, num := range langs {
		res := (float64(num) / float64(total_size)) * 100.0
		results = append(results, res)
		qqq[res] = lang
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(results)))
	for _, percent := range results {
		fmt.Printf(fmtstr, qqq[percent], percent)
	}
	fmt.Printf("---\n%d languages detected in %d bytes of %d files\n", len(langs), total_size, num_files)
}

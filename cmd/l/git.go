package main

import (
	"log"
	"path/filepath"

	"github.com/generaltso/git4go"
	"github.com/generaltso/linguist"
)

func processTree(repo *git4go.Repository, odb *git4go.Odb, tree_id *git4go.Oid, parent []string) {
	var tree *git4go.Tree
	var commit *git4go.Commit
	commit, err := repo.LookupCommit(tree_id)
	if err != nil {
		obj, errr := repo.Lookup(tree_id)
		checkErr(errr)
		switch obj.Type() {
		case git4go.ObjectTree:
			tree = obj.(*git4go.Tree)
		case git4go.ObjectCommit:
			commit = obj.(*git4go.Commit)
		default:
			log.Panicf("%#v not a tree object", obj)
		}
	}
	if commit != nil {
		tree, err = commit.Tree()
		checkErr(err)
	}
	for _, entry := range tree.Entries {
		//fmode := fmt.Sprintf("%06o", int(entry.Filemode))
		ftype := entry.Type.String()
		fhash := entry.Id.String()
		fname := entry.Name

		switch ftype {
		case "tree":
			log.Println("entering subtree", fname)
			oid, err := git4go.NewOid(fhash)
			checkErr(err)
			processTree(repo, odb, oid, append(parent, fname))
		case "blob":
			fname = filepath.Join(append(parent, fname)...)

			oid, err := git4go.NewOid(fhash)
			checkErr(err)
			obj, err := odb.Read(oid)
			checkErr(err)

			size := len(obj.Data)

			log.Println(fname, "is", size, "bytes")
			if size == 0 {
				log.Println(fname, "is empty file, skipping")
				continue
			}

			if !unignore_filenames && linguist.ShouldIgnoreFilename(fname) {
				log.Println(fname, ": filename should be ignored, skipping")
				ignored_paths++
				continue
			}

			by_name := linguist.LanguageByFilename(fname)
			if by_name != "" {
				log.Println(fname, "got result by name: ", by_name)
				putResult(by_name, size)
				continue
			}

			contents := obj.Data

			if !unignore_contents && linguist.ShouldIgnoreContents(contents) {
				log.Println(fname, ": contents should be ignored, skipping")
				ignored_paths++
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

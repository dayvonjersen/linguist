package main

import (
	"io"
	"log"

	"github.com/dayvonjersen/linguist"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func processRepoTreeAt(path string, refName string) {
	log.Println("opening repo at", path)
	repo, err := git.PlainOpen(path)
	checkErr(err)

	log.Println("resolving revision", refName)
	hash, err := repo.ResolveRevision(plumbing.Revision(refName))
	checkErr(err)

	log.Println("getting commit for hash", *hash)
	commit, err := repo.CommitObject(*hash)
	checkErr(err)

	log.Println("getting tree for commit")
	tree, err := commit.Tree()
	checkErr(err)

	seen := make(map[plumbing.Hash]bool)
	walker := object.NewTreeWalker(tree, true, seen)
	defer walker.Close()

	for {
		name, entry, err := walker.Next()
		if err == io.EOF {
			break
		}
		checkErr(err)

		if entry.Mode != filemode.Regular && entry.Mode != filemode.Executable {
			continue
		}

		log.Println("processing", entry.Mode, name)

		blob, err := repo.BlobObject(entry.Hash)
		checkErr(err)

		log.Println(name, "is", blob.Size, "bytes")
		if blob.Size == 0 {
			log.Println(name, "is empty file, skipping")
			continue
		}

		if !unignore_filenames && linguist.ShouldIgnoreFilename(name) {
			log.Println(name, ": filename should be ignored, skipping")
			ignored_paths++
			continue
		}

		langByName := linguist.LanguageByFilename(name)
		if langByName != "" {
			log.Println(name, "got result by name: ", langByName)
			putResult(langByName, blob.Size)
			continue
		}

		r, err := blob.Reader()
		checkErr(err)

		contents, err := io.ReadAll(r)
		checkErr(err)

		if !unignore_contents && linguist.ShouldIgnoreContents(contents) {
			log.Println(name, ": contents should be ignored, skipping")
			ignored_paths++
			continue
		}

		hints := linguist.LanguageHints(name)
		log.Printf("%s got language hints: %#v\n", name, hints)
		langByContents := linguist.LanguageByContents(contents, hints)

		if langByContents != "" {
			log.Println(name, "got result by data: ", langByContents)
			putResult(langByContents, blob.Size)
			continue
		}

		log.Println(name, "got no result!!")
		putResult("(unknown)", blob.Size)
	}
}

package main

import (
	"fmt"

	"github.com/asi-isa/trav"
)

func main() {
	// cla
	const root = "/"
	const threshold = 100_000_000

	where := func(entry trav.FileEntry) bool {
		return entry.Size > threshold
	}

	ch := trav.New().Traverse(root, where)

	report(ch)
}

func report(ch <-chan trav.FileEntry) {
	files := make(map[string][]trav.FileEntry)
	for fileE := range ch {
		files[fileE.Name] = append(files[fileE.Name], fileE)
	}

	duplicateFiles := make(map[string][]trav.FileEntry)
	for key, val := range files {
		if len(val) > 1 {
			duplicateFiles[key] = val
		}
	}

	for fName, files := range duplicateFiles {
		fmt.Println("##################################")
		fmt.Printf("File: %s \n", fName)

		for _, file := range files {
			fmt.Printf("\tPath: %v\n", file.Path)
			fmt.Printf("\tSize: %v\n", file.Size)
			fmt.Printf("\tModified: %v\n\n", file.Date)
		}
		fmt.Println("##################################")
	}
}

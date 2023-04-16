package main

import (
	"fmt"
	"io/fs"
	"log"
	"time"

	"github.com/asi-isa/trav"
)

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type FileEntry struct {
	name     string
	path     string
	size     int64
	modified time.Time
}

func createEntry(path string, entry fs.DirEntry) (FileEntry, bool) {
	info, err := entry.Info()
	handleErr(err)

	size := info.Size()

	const threshold = 100_000_000

	if size > threshold {
		name := info.Name()
		path += "/" + name
		time := info.ModTime()

		return FileEntry{name, path, size, time}, true
	}

	return *new(FileEntry), false
}

func main() {
	// cla
	const root = "/"

	var trav = trav.New[FileEntry](root)
	ch := trav.Traverse(createEntry)

	files := make(map[string][]FileEntry)
	for fileE := range ch {
		files[fileE.name] = append(files[fileE.name], fileE)
	}

	duplicateFiles := make(map[string][]FileEntry)
	for key, val := range files {
		if len(val) > 1 {
			duplicateFiles[key] = val
		}
	}

	for fName, files := range duplicateFiles {
		fmt.Println("##################################")
		fmt.Printf("File: %s \n", fName)

		for _, file := range files {
			fmt.Printf("\tPath: %v\n", file.path)
			fmt.Printf("\tSize: %v\n", file.size)
			fmt.Printf("\tModified: %v\n\n", file.modified)
		}
		fmt.Println("##################################")
	}
}

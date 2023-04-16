package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"sync"
	"time"
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

func traverse(wg *sync.WaitGroup, path string, clb func(string, fs.DirEntry)) {
	dirEntries, err := os.ReadDir(path)

	if err == nil {
		for _, entry := range dirEntries {
			if entry.IsDir() {
				traverse(wg, path+"/"+entry.Name(), clb)
			} else {
				wg.Add(1)
				go clb(path, entry)
			}
		}
	}
}

func main() {
	// cla
	const threshold = 100_000_000
	const path = "/"

	ch := make(chan FileEntry)
	var wg sync.WaitGroup

	sendEntry := func(path string, entry fs.DirEntry) {
		defer wg.Done()

		info, err := entry.Info()
		handleErr(err)

		size := info.Size()

		if size > threshold {

			name := info.Name()
			path += "/" + name
			time := info.ModTime()

			ch <- FileEntry{name, path, size, time}
		}
	}

	go func() {
		defer close(ch)
		defer wg.Wait()

		traverse(&wg, path, sendEntry)
	}()

	files := make(map[string][]FileEntry)
	duplicateFiles := make(map[string][]FileEntry)

	for fileE := range ch {
		files[fileE.name] = append(files[fileE.name], fileE)
	}

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

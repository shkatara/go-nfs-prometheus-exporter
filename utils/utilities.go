package utils

import (
	"log"
	"os"
	"path/filepath"
)

func errorcheck(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func FindDir(path string) []string {
	directories := []string{}
	entries, err := os.ReadDir(path)
	errorcheck(err)
	for _, entry := range entries {
		if (entry.IsDir()) && (entry.Name() != ".git") {
			directories = append(directories, entry.Name())
		}
	}
	return directories
}

func DirSize(dirpath string) int64 {
	var totalSize int64
	err := filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})
	errorcheck(err)
	return totalSize
}

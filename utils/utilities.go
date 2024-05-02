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

func DirSize(target string, dirpath string) float64 {
	var totalSize float64
	newDirPath := target + "/" + dirpath
	err := filepath.Walk(newDirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			totalSize += float64(info.Size())
		}
		return nil
	})
	errorcheck(err)
	return totalSize
}

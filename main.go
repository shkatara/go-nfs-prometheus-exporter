package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func errorcheck(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func lsDir(path string) []string {
	directories := []string{}
	entries, err := os.ReadDir(path)
	errorcheck(err)
	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, entry.Name())
		}
	}
	return directories
}

func main() {
	dir := lsDir("./")
	for _, data := range dir {
		run, _ := exec.Command("du", "-xh", data).Output()
		//output := string(run[3:4])
		output := string(run[:4])
		fmt.Println(output)
	}
}

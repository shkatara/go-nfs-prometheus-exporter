package main

import (
	"fmt"
	"time"

	"example.com/nfs-exporter/utils"
)

func main() {
	var target string = "./"
	dir := utils.FindDir(target)
	for {
		time.Sleep(5 * time.Second)
		for _, data := range dir {
			size := utils.DirSize(data)
			fmt.Println(size)
		}
	}
}

package main

import (
	"os"

	"github.com/0xBB2b/gotree"
)

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	gotree.Print(path)
}

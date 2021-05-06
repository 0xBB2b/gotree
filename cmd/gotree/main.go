package main

import (
	"os"

	"github.com/ne2blink/gotree"
)

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	gotree.Print(path)
}

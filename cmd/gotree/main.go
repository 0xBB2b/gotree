package main

import (
	"os"
	"strings"

	"github.com/ne2blink/gotree"
)

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	path = strings.TrimRight(path, "/")
	path = strings.TrimRight(path, "\\")
	gotree.Print(path)
}

package main

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	tree "github.com/need-being/go-tree"
)

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	root := tree.New(path)
	printer := tree.NewPrinter(os.Stdout, func(n *tree.Node) string {
		if n.Virtual {
			return fmt.Sprint(color.New(color.FgCyan).SprintFunc()(n.Value))
		}
		return fmt.Sprint(n.Value)
	})
	reader(path, root)
	printer.Print(root)
}

func reader(path string, node *tree.Node) error {
	if strings.HasSuffix(path, ".tar") {
		return readTar(path, node)
	}
	return readDir(path, node)
}

func readDir(path string, node *tree.Node) error {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, fi := range dir {
		if fi.IsDir() {
			newNode := node.Add(fi.Name())
			reader(path+"/"+fi.Name()+"/", newNode)
		} else {
			newNode := node.Add(fi.Name())
			if strings.HasSuffix(fi.Name(), ".tar") {
				reader(path+"/"+fi.Name(), newNode)
			}
		}
	}
	return nil
}

func readTar(path string, node *tree.Node) error {
	reader, err := os.Open(path)
	if err != nil {
		return err
	}
	defer reader.Close()

	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		headerName := strings.TrimRight(header.Name, "/")
		node.AddPathString(headerName)
	}
	return nil
}

package gotree

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
)

type FileType uint

const (
	Folder FileType = iota + 1
	File
	Tar
)

const (
	folderColor = color.FgCyan
	fileColor   = color.FgWhite
	tarColor    = color.FgMagenta
)

type URI struct {
	lv   int
	typ  FileType
	name string
}

func Print(path string) {
	if err := printDirs(path); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printDirs(path string) error {
	list, err := readDir(path, 0)
	if err != nil {
		return err
	}
	var tree []string
	for index, uri := range list {
		tab := ""
		isLast := true
		for i := index + 1; i < len(list); i++ {
			if list[i].lv < list[index].lv {
				break
			}
			if list[i].lv == list[index].lv {
				isLast = false
				break
			}
		}
		if index > 0 && len(tree) > 0 && list[index-1].lv > list[index].lv {
			tree = tree[:uri.lv]
		}
		if isLast && uri.typ != File {
			if len(tree) == uri.lv {
				tree = append(tree, "  ")
			}
		} else if uri.typ != File {
			if len(tree) == uri.lv {
				tree = append(tree, "┃ ")
			}
		}
		for i := 0; i < uri.lv; i++ {
			tab = tab + tree[i]
		}
		if index == len(list)-1 || list[index].lv > list[index+1].lv || isLast {
			tab = tab + "┗━"
		} else {
			tab = tab + "┣━"
		}
		fmt.Println(tab + print(uri.typ)(uri.name))
	}
	return nil
}

func readDir(path string, lv int) ([]URI, error) {
	var list []URI
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return list, err
	}
	for _, fi := range dir {
		if fi.IsDir() {
			list = append(list, URI{lv, Folder, fi.Name()})
			l, err := readDir(path+"/"+fi.Name(), lv+1)
			if err != nil {
				return list, err
			}
			list = append(list, l...)
		} else if strings.HasSuffix(fi.Name(), ".tar") {
			list = append(list, URI{lv, Tar, fi.Name()})
			tars, err := readTar(path+"/"+fi.Name(), lv)
			if err != nil {
				return list, err
			}
			list = append(list, tars...)
		} else {
			list = append(list, URI{lv, File, fi.Name()})
		}
	}
	return list, err
}

func readTar(path string, lv int) (list []URI, err error) {
	reader, err := os.Open(path)
	if err != nil {
		return
	}
	defer reader.Close()

	tarReader := tar.NewReader(reader)
	for {
		header, tarErr := tarReader.Next()
		if tarErr == io.EOF {
			break
		} else if tarErr != nil {
			err = tarErr
			return
		}

		fi := header.FileInfo()
		tp := File
		l := len(strings.Split(header.Name, "/")) + lv
		if strings.HasSuffix(header.Name, "/") {
			tp = Folder
			l = l - 1
		}
		list = append(list, URI{l, tp, fi.Name()})
	}
	return
}

func print(filrType FileType) func(a ...interface{}) string {
	switch filrType {
	case Folder:
		return color.New(folderColor).SprintFunc()
	case Tar:
		return color.New(tarColor).SprintFunc()
	default:
		return color.New(fileColor).SprintFunc()
	}
}

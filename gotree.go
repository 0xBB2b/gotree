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
	GeneralFile
	TarFile
)

const (
	tab0        = "  "
	tab1        = "┃ "
	tab2        = "┗━"
	tab3        = "┣━"
	folderColor = color.FgCyan
	fileColor   = color.FgWhite
	tarColor    = color.FgMagenta
)

func Print(path string) {
	if err := printDirs(path, ""); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printDirs(path string, tab string) error {
	dirs, files, err := readDir(path)
	if err != nil {
		return err
	}
	for index, dir := range dirs {
		if index == len(dirs)-1 && len(files) == 0 {
			fmt.Println(tab + tab2 + print(Folder)(dir))
			nextPath := path + "/" + dir
			if err := printDirs(nextPath, tab+tab0); err != nil {
				return err
			}
		} else {
			fmt.Println(tab + tab3 + print(Folder)(dir))
			nextPath := path + "/" + dir
			if err := printDirs(nextPath, tab+tab1); err != nil {
				return err
			}
		}
	}
	for index, file := range files {
		if strings.HasSuffix(file, ".tar") {
			if index == len(files)-1 {
				fmt.Println(tab + tab2 + print(TarFile)(file))
				nextPath := path + "/" + file
				if err := printTar(nextPath, tab+tab0); err != nil {
					return err
				}
			} else {
				fmt.Println(tab + tab3 + print(TarFile)(file))
				nextPath := path + "/" + file
				if err := printTar(nextPath, tab+tab1); err != nil {
					return err
				}
			}
		} else {
			if index == len(files)-1 {
				fmt.Println(tab + tab2 + print(GeneralFile)(file))
			} else {
				fmt.Println(tab + tab3 + print(GeneralFile)(file))
			}
		}
	}

	return nil
}

func printTar(path string, tab string) error {
	reader, err := os.Open(path)
	if err != nil {
		return err
	}
	defer reader.Close()

	tarReader := tar.NewReader(reader)
	var headerNames []string
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		headerNames = append(headerNames, header.Name)
	}
	for index, headerName := range headerNames {
		if index == len(headerNames)-1 {
			fmt.Println(tab + tab2 + print(GeneralFile)(headerName))
		} else {
			fmt.Println(tab + tab3 + print(GeneralFile)(headerName))
		}
	}
	return nil
}

// 读取目录
func readDir(path string) (dirs, flies []string, err error) {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, fi := range dir {
		if fi.IsDir() {
			dirs = append(dirs, fi.Name())
		} else {
			flies = append(flies, fi.Name())
		}
	}
	return
}

// 打印对应的颜色
func print(filrType FileType) func(a ...interface{}) string {
	switch filrType {
	case Folder:
		return color.New(folderColor).SprintFunc()
	case TarFile:
		return color.New(tarColor).SprintFunc()
	default:
		return color.New(fileColor).SprintFunc()
	}
}

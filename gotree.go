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
	folderColor = color.FgCyan
	fileColor   = color.FgWhite
	tarColor    = color.FgMagenta
)

type FileTree struct {
	level    int
	fileType FileType
	name     string
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
	var tabArray []string
	for index, uri := range list {
		tab := ""
		isLast := true
		// 判断是否为当前文件将中最后一个
		for i := index + 1; i < len(list); i++ {
			if list[i].level < list[index].level {
				break
			}
			if list[i].level == list[index].level {
				isLast = false
				break
			}
		}
		// 删除多余的制表符
		if index > 0 && len(tabArray) > 0 && list[index-1].level > list[index].level {
			tabArray = tabArray[:uri.level]
		}
		// 添加对应的制表符
		if isLast && uri.fileType != GeneralFile {
			if len(tabArray) == uri.level {
				tabArray = append(tabArray, "  ")
			}
		} else if uri.fileType != GeneralFile {
			if len(tabArray) == uri.level {
				tabArray = append(tabArray, "┃ ")
			}
		}
		for i := 0; i < uri.level; i++ {
			tab = tab + tabArray[i]
		}
		if index == len(list)-1 || list[index].level > list[index+1].level || isLast {
			tab = tab + "┗━"
		} else {
			tab = tab + "┣━"
		}
		fmt.Println(tab + print(uri.fileType)(uri.name))
	}
	return nil
}

// 读取目录
func readDir(path string, level int) ([]FileTree, error) {
	var list []FileTree
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return list, err
	}
	for _, fi := range dir {
		if fi.IsDir() {
			list = append(list, FileTree{level, Folder, fi.Name()})
			l, err := readDir(path+"/"+fi.Name(), level+1)
			if err != nil {
				return list, err
			}
			list = append(list, l...)
		} else if strings.HasSuffix(fi.Name(), ".tar") {
			list = append(list, FileTree{level, TarFile, fi.Name()})
			tars, err := readTar(path+"/"+fi.Name(), level)
			if err != nil {
				return list, err
			}
			list = append(list, tars...)
		} else {
			list = append(list, FileTree{level, GeneralFile, fi.Name()})
		}
	}
	return list, err
}

// 读取tar文件
func readTar(path string, level int) (list []FileTree, err error) {
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

		fileInfo := header.FileInfo()
		fileType := GeneralFile
		lv := len(strings.Split(header.Name, "/")) + level
		if strings.HasSuffix(header.Name, "/") {
			fileType = Folder
			lv = lv - 1
		}
		list = append(list, FileTree{lv, fileType, fileInfo.Name()})
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

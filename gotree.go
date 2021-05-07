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
	level    int      // 层级
	fileType FileType // 类型
	name     string   // 名称
}

func Print(path string) {
	if err := printDirs(path); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printDirs(path string) error {
	files, err := readDir(path, 0)
	if err != nil {
		return err
	}
	var tabArray []string
	for index, file := range files {
		tab := ""
		isLast := true
		// 判断是否为当前文件夹中最后一个
		for i := index + 1; i < len(files); i++ {
			if files[i].level < files[index].level {
				break
			}
			if files[i].level == files[index].level {
				isLast = false
				break
			}
		}
		// 删除多余的制表符
		if index > 0 && len(tabArray) > 0 && files[index-1].level > files[index].level {
			tabArray = tabArray[:file.level]
		}
		// 添加对应的制表符
		if isLast && file.fileType != GeneralFile {
			if len(tabArray) == file.level {
				tabArray = append(tabArray, "  ")
			}
		} else if file.fileType != GeneralFile {
			if len(tabArray) == file.level {
				tabArray = append(tabArray, "┃ ")
			}
		}
		for i := 0; i < file.level; i++ {
			tab = tab + tabArray[i]
		}
		if index == len(files)-1 || files[index].level > files[index+1].level || isLast {
			tab = tab + "┗━"
		} else {
			tab = tab + "┣━"
		}
		fmt.Println(tab + print(file.fileType)(file.name))
	}
	return nil
}

// 读取目录
func readDir(path string, level int) ([]FileTree, error) {
	var files []FileTree
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return files, err
	}
	for _, fi := range dir {
		if fi.IsDir() {
			files = append(files, FileTree{level, Folder, fi.Name()})
			fs, err := readDir(path+"/"+fi.Name(), level+1)
			if err != nil {
				return files, err
			}
			files = append(files, fs...)
		} else if strings.HasSuffix(fi.Name(), ".tar") {
			files = append(files, FileTree{level, TarFile, fi.Name()})
			tars, err := readTar(path+"/"+fi.Name(), level)
			if err != nil {
				return files, err
			}
			files = append(files, tars...)
		} else {
			files = append(files, FileTree{level, GeneralFile, fi.Name()})
		}
	}
	return files, err
}

// 读取tar文件
func readTar(path string, level int) ([]FileTree, error) {
	var files []FileTree
	reader, err := os.Open(path)
	if err != nil {
		return files, err
	}
	defer reader.Close()

	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return files, err
		}

		fileInfo := header.FileInfo()
		fileType := GeneralFile
		lv := len(strings.Split(header.Name, "/")) + level
		if strings.HasSuffix(header.Name, "/") {
			fileType = Folder
			lv = lv - 1
		}
		files = append(files, FileTree{lv, fileType, fileInfo.Name()})
	}
	return files, err
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

package gotree

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
)

func Print(path string) {
	if err := printDirs(path, ""); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printDirs(path string, n string) error {
	dirs, files, err := readDir(path)
	if err != nil {
		return err
	}
	for i, dir := range dirs {
		if len(dirs)-1 == i && len(files) == 0 {
			fmt.Println(n+"┗━", color.CyanString(dir))
			if err := printDirs(path+"/"+dir+"/", n+"   "); err != nil {
				return err
			}
		} else {
			fmt.Println(n+"┣━", color.CyanString(dir))
			if err := printDirs(path+"/"+dir+"/", n+"┃  "); err != nil {
				return err
			}
		}
	}
	for i, file := range files {
		if len(files)-1 == i {
			color.White(n + "┗━ " + file)
		} else {
			color.White(n + "┣━ " + file)
		}
	}
	return nil
}

func readDir(path string) (dirs []string, files []string, err error) {
	dirs = make([]string, 0)
	files = make([]string, 0)
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, fi := range dir {
		if fi.IsDir() {
			dirs = append(dirs, fi.Name())
		} else {
			files = append(files, fi.Name())
		}
	}
	return
}

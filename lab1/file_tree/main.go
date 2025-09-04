package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func subTree(out io.Writer, path string, printFiles bool, prefix string) error {
	dir, err := os.ReadDir(path)
	if err != nil {
		return errors.New("Ошибка")
	}

	if !printFiles {
		var filteredDir []os.DirEntry
		for _, entry := range dir {
			if entry.IsDir() {
				filteredDir = append(filteredDir, entry)
			}
		}
		dir = filteredDir
	}

	sort.Slice(dir, func(i, j int) bool {
		return dir[i].Name() < dir[j].Name()
	})

	for i, file := range dir {
		connector := "├───"
		newPrefix := prefix + "│\t"
		if i == len(dir)-1 {
			connector = "└───"
			newPrefix = prefix + "\t"
		}

		info, _ := file.Info()
		fileSize := info.Size()
		size := strconv.Itoa(int(fileSize))

		var sizeStr string
		if size == "0" {
			sizeStr = " (empty)"
		} else {
			sizeStr = " (" + size + "b)"
		}

		if file.IsDir() {
			newPath := filepath.Join(path, file.Name())
			fmt.Fprintln(out, prefix+connector+file.Name())
			err := subTree(out, newPath, printFiles, newPrefix)
			if err != nil {
				return errors.New("Ошибка")
			}
		} else if printFiles {
			fmt.Fprintln(out, prefix+connector+file.Name()+sizeStr)
		}
	}
	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return subTree(out, path, printFiles, "")
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

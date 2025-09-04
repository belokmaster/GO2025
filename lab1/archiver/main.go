package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func addFileToZip(zipWriter *zip.Writer, fullPath string, shortPath string) error {
	openOldFile, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	defer openOldFile.Close()

	writer, err := zipWriter.Create(shortPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, openOldFile)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if !(len(os.Args) == 3) {
		panic("need 3 args command line")
	}

	oldFile := os.Args[1]
	newFile := os.Args[2]

	zipFile, err := os.Create(newFile)
	if err != nil {
		panic("problem to create zip archive")
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	fileInfo, err := os.Stat(oldFile)
	if err != nil {
		panic("problem to get info about oldfile")
	}

	if fileInfo.IsDir() {
		err := filepath.Walk(oldFile, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			shortPath, err := filepath.Rel(oldFile, path)
			return addFileToZip(zipWriter, path, shortPath)
		})

		if err != nil {
			panic(err)
		}
	} else {
		shortPath := filepath.Base(oldFile)
		err = addFileToZip(zipWriter, oldFile, shortPath)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("File successfully archived")
}

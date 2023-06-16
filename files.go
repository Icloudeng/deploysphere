package main

import (
	"errors"
	"log"
	"os"
)

type file struct{}

var File = file{}

func (f *file) checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}

func (f *file) createIfNotExists(filePath string) bool {
	isFileExist := f.checkFileExists(filePath)

	if !isFileExist {
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	return isFileExist
}

func (f *file) writeInFile(filePath string, content string) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		log.Fatal(err)
	}
}

func (f *file) readFile(filePath string) []byte {
	content, err := os.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	return content
}

func (f *file) createIfNotExistsWithContent(filePath string, content string) {
	isFileExist := f.createIfNotExists(filePath)

	if !isFileExist {
		f.writeInFile(filePath, content)
	}
}

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

func (f *file) createIfNotExists(filePath string) {
	isFileExist := f.checkFileExists(filePath)

	if !isFileExist {
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}
}

func (f *file) writeInFile(filePath string, content string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		log.Fatal(err)
	}
}

func (f *file) createIfNotExistsWithContent(filePath string, content string) {
	f.createIfNotExists(filePath)
	f.writeInFile(filePath, content)
}

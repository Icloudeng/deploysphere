package files

import (
	"errors"
	"log"
	"os"
)

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}

func createIfNotExists(filePath string) bool {
	isFileExist := checkFileExists(filePath)

	if !isFileExist {
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	return isFileExist
}

func WriteInFile(filePath string, content string) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		log.Fatal(err)
	}
}

func ReadFile(filePath string) []byte {
	content, err := os.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	return content
}

func CreateIfNotExistsWithContent(filePath string, content string) {
	isFileExist := createIfNotExists(filePath)

	if !isFileExist {
		WriteInFile(filePath, content)
	}
}

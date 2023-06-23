package files

import (
	"errors"
	"log"
	"os"
	"path"
)

func GetPwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot the current dir %s", err)
	}

	return dir
}

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

func ExistsProvisionerPlaformReadDir(platform string) bool {
	pwd := GetPwd()

	entries, err := os.ReadDir(path.Join(pwd, "infrastrure/provisioner/scripts/platforms"))
	if err != nil {
		log.Panicf("failed reading directory: %s", err)
	}

	exists := false
	for _, v := range entries {
		if v.Name() == platform && v.IsDir() {
			exists = true
		}
	}

	return exists
}

func ReadPlatformMetadataFields() []byte {
	pwd := GetPwd()
	return ReadFile(path.Join(pwd, "infrastrure/provisioner/scripts/metadata-required.json"))
}

func ReadProvisionerPlaforms() []string {
	pwd := GetPwd()
	entries, err := os.ReadDir(path.Join(pwd, "infrastrure/provisioner/scripts/platforms"))
	if err != nil {
		log.Panicf("failed reading directory: %s", err)
	}

	var platforms []string

	for _, v := range entries {
		if v.IsDir() {
			platforms = append(platforms, v.Name())
		}
	}

	return platforms
}

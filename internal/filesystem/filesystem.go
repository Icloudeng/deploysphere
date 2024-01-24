package filesystem

import (
	"errors"
	"log"
	"os"
	"path"
)

var ProvisionerDir string
var TerraformDir string

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
			log.Fatalln(err)
		}
		defer file.Close()
	}

	return isFileExist
}

func WriteInFile(filePath string, content string) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		log.Fatalln(err)
	}
}

func ReadFile(filePath string) []byte {
	content, err := os.ReadFile(filePath)

	if err != nil {
		log.Fatalln(err)
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
	entries, err := os.ReadDir(path.Join(ProvisionerDir, "scripts/platforms"))
	if err != nil {
		log.Panicln("failed reading directory:", err)
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
	return ReadFile(path.Join(ProvisionerDir, "scripts/platform-meta-fields.json"))
}

func ReadEnvironmentsFile() []byte {
	return ReadFile(path.Join(ProvisionerDir, "scripts/environments.json"))
}

func ReadConfigurationMetadataFields() []byte {
	return ReadFile(path.Join(ProvisionerDir, "scripts/platform-configuration-fields.json"))
}

func ReadProvisionerPlatforms() []string {
	entries, err := os.ReadDir(path.Join(ProvisionerDir, "scripts/platforms"))
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

func init() {
	pwd := GetPwd()
	ProvisionerDir = path.Join(pwd, "infrastructure/provisioner")
	TerraformDir = path.Join(pwd, "infrastructure/terraform")
}

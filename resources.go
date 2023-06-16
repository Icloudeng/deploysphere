package main

import (
	"log"
	"os"
	"path"
)

type terraformResourceFiles struct {
	Ovh     string
	Proxmox string
}

var TfResource = &terraformResourceFiles{}

func (f *terraformResourceFiles) GetFiles() [2]string {
	return [...]string{f.Ovh, f.Proxmox}
}

func initTerraformResourceFiles() {
	file := "resource_auto.tf.json"

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot get the current dir %s", err)
	}

	// Store resource paths
	TfResource.Ovh = path.Join(pwd, "infrastrure/terraform/modules/ovh", file)
	TfResource.Proxmox = path.Join(pwd, "infrastrure/terraform/modules/proxmox", file)

	// Create file if not exist
	for _, file := range TfResource.GetFiles() {
		File.createIfNotExistsWithContent(file, "{}")
	}
}

func init() {
	// Create file id not exist yet
	initTerraformResourceFiles()
}

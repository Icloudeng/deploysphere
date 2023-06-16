package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

type terrafrom struct {
	tk *tfexec.Terraform
}

var Tf = &terrafrom{}

func init() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot the current dir %s", err)
	}

	installer := &releases.ExactVersion{
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion("1.5.0")),
		InstallDir: path.Join(dir, "./bin"),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	workingDir := path.Join(dir, "infrastrure/terraform")
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		log.Fatalf("error running Show: %s", err)
	}

	fmt.Printf("Terraform initialized %v", &state)

	Tf.tk = tf
}

func (t *terrafrom) plan() {
	tf := t.tk
	ctx := context.Background()

	options := []tfexec.PlanOption{
		tfexec.VarFile("variables.tfvars"),
	}

	state, err := tf.Plan(ctx, options...)
	if err != nil {
		log.Fatalf("error running Show: %s", err)
	}

	log.Printf("Terraform plan state: %v", state)
}

func (t *terrafrom) apply() {
	tf := t.tk
	ctx := context.Background()

	options := []tfexec.ApplyOption{
		tfexec.VarFile("variables.tfvars"),
	}

	err := tf.Apply(ctx, options...)
	if err != nil {
		log.Fatalf("Error running Show: %s", err)
	}

	log.Printf("********* Terraform applied ! ***********")
}

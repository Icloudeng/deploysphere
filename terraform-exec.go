package main

import (
	"context"
	"log"
	"path"
	"smatflow/platform-installer/files"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

type terrafrom struct {
	tk *tfexec.Terraform
}

var Tf = &terrafrom{}

func init() {
	pwd := files.GetPwd()

	installer := &releases.ExactVersion{
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion("1.5.0")),
		InstallDir: path.Join(pwd, "./bin"),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	workingDir := path.Join(pwd, "infrastrure/terraform")

	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

	Tf.tk = tf
}

func (t *terrafrom) Plan() {
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

func (t *terrafrom) Show() *tfjson.StateModule {
	tf := t.tk
	ctx := context.Background()

	state, err := tf.Show(ctx)

	if err != nil || state == nil || state.Values == nil {
		return nil
	}

	return state.Values.RootModule
}

func (t *terrafrom) Apply() {
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

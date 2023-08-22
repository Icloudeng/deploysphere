package terraform

import (
	"context"
	"log"
	"path"
	"smatflow/platform-installer/pkg/filesystem"
	"smatflow/platform-installer/pkg/pubsub"
	"smatflow/platform-installer/pkg/structs"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

type exec struct {
	tk *tfexec.Terraform
}

var Exec = &exec{}

func init() {
	pwd := filesystem.GetPwd()

	installer := &releases.ExactVersion{
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion("1.5.4")),
		InstallDir: path.Join(pwd, "./bin"),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("Error installing Terraform: %s", err)
	}

	tf, err := tfexec.NewTerraform(filesystem.TerraformDir, execPath)
	if err != nil {
		log.Fatalf("Error running NewTerraform: %s", err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("Error running Init: %s", err)
	}

	Exec.tk = tf
}

func (t *exec) Plan(notifier bool) error {
	options := []tfexec.PlanOption{
		tfexec.VarFile("variables.tfvars"),
	}

	if state, err := t.tk.Plan(context.Background(), options...); err != nil {
		if notifier {
			go pubsub.BusEvent.Publish(pubsub.RESOURCES_NOTIFIER_EVENT, structs.Notifier{
				Status:  "failed",
				Details: "Terraform Plan",
				Logs:    err.Error(),
			})
		}

		log.Printf("error running Show: %s", err.Error())
		return err
	} else {
		log.Printf("Terraform plan state: %v", state)
	}

	return nil
}

func (t *exec) Show() *tfjson.StateModule {
	state, err := t.tk.Show(context.Background())

	if err != nil || state == nil || state.Values == nil {
		return nil
	}

	return state.Values.RootModule
}

func (t *exec) Apply(notifier bool) error {
	if err := t.Plan(notifier); err != nil {
		return err
	}

	options := []tfexec.ApplyOption{
		tfexec.VarFile("variables.tfvars"),
	}

	if err := t.tk.Apply(context.Background(), options...); err != nil {
		if notifier {
			go pubsub.BusEvent.Publish(pubsub.RESOURCES_NOTIFIER_EVENT, structs.Notifier{
				Status:  "failed",
				Details: "Terraform Apply",
				Logs:    err.Error(),
			})
		}

		log.Printf("Error running Show: %s", err)
		return err
	}

	log.Printf("********* Terraform applied ! ***********")

	return nil
}

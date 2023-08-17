package terraform

import (
	"context"
	"log"
	"path"
	"smatflow/platform-installer/pkg/files"
	"smatflow/platform-installer/pkg/pubsub"
	"smatflow/platform-installer/pkg/structs"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

type Terrafrom struct {
	tk *tfexec.Terraform
}

var Tf = &Terrafrom{}

func init() {
	pwd := files.GetPwd()

	installer := &releases.ExactVersion{
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion("1.5.4")),
		InstallDir: path.Join(pwd, "./bin"),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	tf, err := tfexec.NewTerraform(files.TerraformDir, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

	Tf.tk = tf
}

func (t *Terrafrom) Plan(notifier bool) error {
	tf := t.tk
	ctx := context.Background()

	options := []tfexec.PlanOption{
		tfexec.VarFile("variables.tfvars"),
	}

	state, err := tf.Plan(ctx, options...)
	if err != nil {

		if notifier {
			go func() {
				pubsub.BusEvent.Publish(pubsub.RESOURCES_NOTIFIER_EVENT, structs.Notifier{
					Status:  "failed",
					Details: "Terraform Plan",
					Logs:    err.Error(),
				})
			}()
		}

		log.Printf("error running Show: %s", err.Error())
		return err
	}

	log.Printf("Terraform plan state: %v", state)

	return nil
}

func (t *Terrafrom) Show() *tfjson.StateModule {
	tf := t.tk
	ctx := context.Background()

	state, err := tf.Show(ctx)

	if err != nil || state == nil || state.Values == nil {
		return nil
	}

	return state.Values.RootModule
}

func (t *Terrafrom) Apply(notifier bool) error {
	if err := t.Plan(notifier); err != nil {
		return err
	}

	options := []tfexec.ApplyOption{
		tfexec.VarFile("variables.tfvars"),
	}

	err := t.tk.Apply(context.Background(), options...)
	if err != nil {
		if notifier {
			go func() {
				pubsub.BusEvent.Publish(pubsub.RESOURCES_NOTIFIER_EVENT, structs.Notifier{
					Status:  "failed",
					Details: "Terraform Apply",
					Logs:    err.Error(),
				})
			}()
		}

		log.Printf("Error running Show: %s", err)
		return err
	}

	log.Printf("********* Terraform applied ! ***********")

	return nil
}

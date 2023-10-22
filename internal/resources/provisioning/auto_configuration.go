package provisioning

import (
	"os"
	"os/exec"

	"github.com/icloudeng/platform-installer/internal/filesystem"
	"github.com/icloudeng/platform-installer/internal/structs"
)

func CreateAutoConfigurationProvisioning(params structs.AutoConfiguration) error {

	cmd := exec.Command(
		"bash", "auto-configuration.sh",
		"--type", params.Type,
		"--platform", params.Platform,
		"--reference", params.PlatformRef,
		"--config-reference", params.PlatformConfigRef,
	)

	cmd.Dir = filesystem.ProvisionerDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

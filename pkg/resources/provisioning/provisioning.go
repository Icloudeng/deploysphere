package provisioning

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"smatflow/platform-installer/pkg/filesystem"
	"smatflow/platform-installer/pkg/structs"
)

func provisioning(prov structs.Provisioning, file string, jobId uint) {
	platform := prov.Platform

	metadata, _ := json.Marshal(platform.Metadata)
	metadatab64 := base64.StdEncoding.EncodeToString(metadata)

	cmd := exec.Command(
		"bash", file,
		"--ansible-user", prov.MachineUser,
		"--vmip", prov.MachineIp,
		"--platform", platform.Name,
		"--metadata", metadatab64,
		"--job-id", fmt.Sprintf("%d", jobId),
		// "--reference", prov.Ref, Don't uncomment this line, can cause mis functioning from redis pubsub
	)

	cmd.Dir = filesystem.ProvisionerDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}

func CreatePlatformProvisioning(prov structs.Provisioning, jobId uint) {
	provisioning(prov, "installer.sh", jobId)
}

func CreateConfigurationProvisioning(prov structs.Provisioning, jobId uint) {
	provisioning(prov, "configuration.sh", jobId)
}

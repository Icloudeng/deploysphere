package provisioning

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"os/exec"
	"smatflow/platform-installer/lib/files"
	"smatflow/platform-installer/lib/structs"
)

func CreateProvisioning(prov structs.Provisioning) {
	platform := prov.Platform

	metadata, _ := json.Marshal(platform.Metadata)
	metadatab64 := base64.StdEncoding.EncodeToString(metadata)

	cmd := exec.Command(
		"bash", "installer.sh",
		"--ansible-user", prov.MachineUser,
		"--vmip", prov.MachineIp,
		"--platform", platform.Name,
		"--metadata", metadatab64,
	)

	cmd.Dir = files.ProvisionerDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}

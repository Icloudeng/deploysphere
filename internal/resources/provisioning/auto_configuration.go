package provisioning

import (
	"os"
	"os/exec"
	"strings"

	"github.com/icloudeng/platform-installer/internal/filesystem"
	"github.com/icloudeng/platform-installer/internal/structs"
)

func CreateAutoConfigurationProvisioning(params structs.AutoConfiguration) ([]byte, error) {

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

	return cmd.CombinedOutput()
}

func ExtractDataFromConfigurationOutputCommand(output []byte) string {
	// Convert the command output to a string
	outputStr := string(output)

	// Split the output into lines
	lines := strings.Split(outputStr, "\n")

	// Extract data between %%...%% markers
	var extractedData string
	for _, line := range lines {
		start := strings.Index(line, "%%")
		end := strings.LastIndex(line, "%%")
		if start != -1 && end != -1 && start < end {
			data := line[start+2 : end]
			extractedData += data
		}
	}

	return extractedData
}

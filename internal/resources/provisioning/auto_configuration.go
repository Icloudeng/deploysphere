package provisioning

import (
	"fmt"
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

	fmt.Printf("CMD string: %s\n", cmd.String())
	fmt.Printf("CMD Args: %s\n", cmd.Args)

	cmd.Dir = filesystem.ProvisionerDir

	output, err := cmd.CombinedOutput()

	if cmd.Process != nil && err != nil {
		cmd.Process.Kill()
	}

	return output, err
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

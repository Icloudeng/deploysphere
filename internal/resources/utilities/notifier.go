package utilities

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"

	"github.com/icloudeng/platform-installer/internal/filesystem"
	"github.com/icloudeng/platform-installer/internal/structs"
)

func SendNotification(notifier structs.Notifier) {

	// Clear up
	cmd := exec.Command(
		"bash", "notifier.sh",
		"--status", notifier.Status,
		"--details", notifier.Details,
		"--logs", base64.StdEncoding.EncodeToString([]byte(notifier.Logs)),
		"--metadata", base64.StdEncoding.EncodeToString([]byte(notifier.Metadata)),
		"--slicetop", "true",
	)

	cmd.Dir = filesystem.ProvisionerDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	fmt.Println("Notification Send !")
}

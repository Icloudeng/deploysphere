package proxyhost

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"smatflow/platform-installer/lib/files"
	"smatflow/platform-installer/lib/structs"
)

func DeleteProxyHost(domain string) {
	// Encode domain
	idomain := map[string]string{"domain": domain}
	metadata, _ := json.Marshal(idomain)
	// Clear up
	cmd := exec.Command(
		"bash", "./nginx-pm.sh",
		"--action delete",
		fmt.Sprintf("--metadata %s", base64.StdEncoding.EncodeToString(metadata)),
	)

	cmd.Dir = files.ProvisionerDir

	if err := cmd.Run(); err != nil {
		log.Panicln(err)
	}

	fmt.Println("Delete Proxy Host, Done!")
}

func CreateProxyHost(proxyhost structs.ProxyHost) {
	// Encode domain
	idomain := map[string]string{"domain": proxyhost.Domain}
	metadata, _ := json.Marshal(idomain)

	// Clear up
	cmd := exec.Command(
		"bash", "nginx-pm.sh",
		"--action create",
		fmt.Sprintf("--platform %s", proxyhost.Platform),
		fmt.Sprintf("--ip %s", proxyhost.Ip),
		fmt.Sprintf("--metadata %s",
			base64.StdEncoding.EncodeToString(metadata),
		),
	)

	cmd.Dir = files.ProvisionerDir

	cmd.Run()

	fmt.Println("Create Proxy Host, Done!")
}

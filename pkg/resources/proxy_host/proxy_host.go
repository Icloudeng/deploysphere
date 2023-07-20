package proxyhost

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"smatflow/platform-installer/pkg/files"
	"smatflow/platform-installer/pkg/structs"
	"strings"
)

func cleanDomain(url string) string {
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimSuffix(url, "/")
	return url
}

func DeleteProxyHost(domain string) {
	// Encode domain
	idomain := map[string]string{"domain": cleanDomain(domain)}
	metadata, _ := json.Marshal(idomain)
	// Clear up
	cmd := exec.Command(
		"bash", "nginx-pm.sh",
		"--action", "delete",
		"--metadata", base64.StdEncoding.EncodeToString(metadata),
	)

	cmd.Dir = files.ProvisionerDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	fmt.Println("Delete Proxy Host, Done!")
}

func CreateProxyHost(proxyhost structs.ProxyHost) {
	// Encode domain
	idomain := map[string]string{"domain": cleanDomain(proxyhost.Domain)}
	metadata, _ := json.Marshal(idomain)

	// Clear up
	cmd := exec.Command(
		"bash", "nginx-pm.sh",
		"--action", "create",
		"--platform", proxyhost.Platform,
		"--ip", proxyhost.Hostname,
		"--metadata", base64.StdEncoding.EncodeToString(metadata),
	)

	cmd.Dir = files.ProvisionerDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	fmt.Println("Create Proxy Host, Done!")
}

package proxmox

import (
	"fmt"
	"smatflow/platform-installer/pkg/env"
	"strconv"
	"strings"

	"github.com/luthermonson/go-proxmox"
)

var PmClient *proxmox.Client

const MIN_CPU = 2
const MIN_DISK = 53_687_091_200

func init() {
	PmClient = proxmox.NewClient(env.Config.PROXMOX_API_URL)
	if err := PmClient.Login(env.Config.PROXMOX_USERNAME, env.Config.PROXMOX_PASSWORD); err != nil {
		panic(err)
	}

	version, err := PmClient.Version()
	if err != nil {
		panic(err)
	}

	fmt.Println("Proxmox Version: " + version.Release)
}

func VmQemuIDExists(id int) bool {
	PmClient.Login(env.Config.PROXMOX_USERNAME, env.Config.PROXMOX_PASSWORD)
	cluster, err := PmClient.Cluster()
	if err != nil {
		return false
	}

	resources, _ := cluster.Resources()

	for _, resource := range resources {
		if resource.Type == "qemu" {
			if s := strings.Split(resource.ID, "/")[1]; s == strconv.Itoa(id) {
				return true
			}
		}
	}

	return false
}

func SelectNodeWithMostResources() (*proxmox.NodeStatus, error) {
	PmClient.Login(env.Config.PROXMOX_USERNAME, env.Config.PROXMOX_PASSWORD)
	nodes, err := PmClient.Nodes()
	if err != nil {
		return nil, err
	}

	var minRam uint64
	var selectedNode *proxmox.NodeStatus

	for _, node := range nodes {
		if ram := node.MaxMem - node.Mem; ram > minRam {
			selectedNode = node
			minRam = ram
		}
	}

	return selectedNode, nil
}

package structs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// Attribute to Ignored changes Terraform Lifecycle
var IGNORE_CHANGES []string = []string{
	"network",
	"disk",
	"ciuser",
	"id",
	"name",
	"qemu_os",
	"sshkeys",
	"ipconfig0",
	"ipconfig1",
	"vmid",
	"clone",
}

type ProxmoxVmQemu struct {
	Name        string `json:"name" binding:"required,resourceref"`
	TargetNode  string `json:"target_node" binding:"required"`
	Vmid        int    `json:"vmid" binding:"omitempty,number"`
	Description string `json:"desc"`
	Clone       string `json:"clone" binding:"required"`

	FullClone bool `json:"full_clone" binding:"boolean"`
	// OS_Network_Config string `json:"os_network_config"`

	// Cloud Init
	OsType    string `json:"os_type" binding:"oneof=ubuntu centos cloud-init"`
	IpConfig0 string `json:"ipconfig0"`
	IpConfig1 string `json:"ipconfig1"`

	OnBoot bool   `json:"onboot"  binding:"boolean"`
	Agent  int    `json:"agent"   binding:"number,oneof=1 0"`
	Memory int    `json:"memory"  binding:"number,min=1024"` //RAM
	Scsihw string `json:"scsihw"  binding:"oneof=lsi lsi53c810 megasas pvscsi virtio-scsi-pci virtio-scsi-single"`

	Cores   int    `json:"cores" binding:"number,min=1,max=10"`   //CPU
	Sockets int    `json:"sockets" binding:"number,min=1,max=10"` //CPU
	Cpu     string `json:"cpu"`                                   //CPU
	Numa    bool   `json:"numa" binding:"boolean"`
	Tags    string `json:"tags"`

	Network   []*PmVmQemuNetwork     `json:"network"`
	Lifecycle []*PmResourceLifecycle `json:"lifecycle"`

	Provisioner [1]interface{} `json:"provisioner"`
}

type PmResourceLifecycle struct {
	IgnoreChanges []string `json:"ignore_changes"`
}

type PmVmQemuNetwork struct {
	Bridge  string `json:"bridge" binding:"required"`
	Model   string `json:"model" binding:"required,oneof=e1000 e1000-82540em e1000-82544gc e1000-82545em i82551 i82557b i82559er ne2k_isa ne2k_pci pcnet rtl8139 virtio vmxnet3"`
	Macaddr string `json:"macaddr"`
	Tag     int    `json:"tag" binding:"number"`
}

// local-exec
type PmLocalExec struct {
	Command    string `json:"command"`
	WorkingDir string `json:"working_dir"`
}

type PmLocalExecProvisioner struct {
	LocalExec [1]*PmLocalExec `json:"local-exec"`
}

// remote-exec
type PmRemoteExec struct {
	Inline *[]string `json:"inline"`
}

type PmRemoteExecProvisioner struct {
	RemoteExec [1]*PmRemoteExec `json:"remote-exec"`
}

type ResetProxmoxVmQemuFields struct {
	Vm       *ProxmoxVmQemu
	Platform Platform
	Ref      string
	JobID    uint
}

func NewProxmoxVmQemu(ref string) *ProxmoxVmQemu {
	vm := ProxmoxVmQemu{
		Vmid:      0,
		FullClone: true,
		OsType:    "cloud-init",
		OnBoot:    true,
		Agent:     1,
		Memory:    2048,
		Scsihw:    "virtio-scsi-pci",
		Cores:     2,
		Sockets:   1,
		Cpu:       "host",
		Numa:      true,
		Tags:      "platform-installer",
		IpConfig0: "ip6=auto,ip=dhcp",
		IpConfig1: "ip6=auto,ip=dhcp",
	}

	vm.Network = append(vm.Network, &PmVmQemuNetwork{
		Bridge: "vmbr0",
		Model:  "virtio",
		Tag:    -1,
	})

	ResetUnmutableProxmoxVmQemu(ResetProxmoxVmQemuFields{
		Vm:       &vm,
		Platform: Platform{},
		Ref:      ref,
	})

	return &vm
}

func newProxmoxResourceLifecycle() *PmResourceLifecycle {
	lifecycle := PmResourceLifecycle{}

	lifecycle.IgnoreChanges = append(
		lifecycle.IgnoreChanges,
		IGNORE_CHANGES...,
	)

	return &lifecycle
}

func newProxmoxProvisioner(platform Platform, ref string, jobid uint) [1]interface{} {
	// Provisioner local-exec
	local_exec := &PmLocalExecProvisioner{}

	if len(platform.Name) > 0 {
		name := platform.Name
		metadata, _ := json.Marshal(platform.Metadata)
		metadatab64 := base64.StdEncoding.EncodeToString(metadata)

		local_exec.LocalExec[0] = &PmLocalExec{
			// Run our ansible scripts here
			Command: fmt.Sprintf("chmod +x installer.sh && ./installer.sh --ansible-user ${self.ciuser} --vmip ${self.default_ipv4_address} --reference %s --job-id %d --platform %s --metadata %s", ref, jobid, name, metadatab64),
			// Relative to infrastructure/terraform
			WorkingDir: "../provisioner",
		}
	} else {
		local_exec.LocalExec[0] = &PmLocalExec{
			// Run our ansible scripts here
			Command: "echo hey....",
			// Relative to infrastructure/terraform
			WorkingDir: "../provisioner",
		}
	}

	// Provisioner remote-exec
	// remote_exec := &PmRemoteExecProvisioner{}
	// remote_exec.RemoteExec[0] = &PmRemoteExec{
	// 	// Sample message to display vm
	// 	Inline: &[]string{"Cool, we are ready for provisioning"},
	// }

	return [1]interface{}{local_exec}
}

func ResetUnmutableProxmoxVmQemu(data ResetProxmoxVmQemuFields) {
	data.Vm.Lifecycle = nil
	data.Vm.Lifecycle = append(data.Vm.Lifecycle, newProxmoxResourceLifecycle())

	data.Vm.Provisioner = newProxmoxProvisioner(data.Platform, data.Ref, data.JobID)
}

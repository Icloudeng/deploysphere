#!/bin/bash
export UBUNTU_VERSION=focal
export TEMPLATE_NAME="ubuntu-20.04-cloudinit-template"
export TEMPLATE_ID=9004
export TEMPLATE_IMAGE="${UBUNTU_VERSION}-server-cloudimg-amd64"
export VM_DISK="local"
export VM_BRIDGE="vmbr1"

sudo apt update -y
sudo apt install libguestfs-tools -y
if [ ! -f "${TEMPLATE_IMAGE}.img" ]; then
    wget https://cloud-images.ubuntu.com/${UBUNTU_VERSION}/current/${TEMPLATE_IMAGE}.img
fi

sudo virt-customize -a ${TEMPLATE_IMAGE}.img --install qemu-guest-agent,net-tools,bash-completion --run-command 'systemctl enable qemu-guest-agent.service' --truncate /etc/machine-id

qemu-img resize ${TEMPLATE_IMAGE}.img 30G
qm create ${TEMPLATE_ID} --name ${TEMPLATE_NAME} --memory 2048 --cores 2 --net0 virtio,bridge=${VM_BRIDGE} --scsihw virtio-scsi-pci
qm set ${TEMPLATE_ID} --scsi0 ${VM_DISK}:0,import-from=/root/${TEMPLATE_IMAGE}.img
qm set ${TEMPLATE_ID} --ide2 ${VM_DISK}:cloudinit
qm set ${TEMPLATE_ID} --boot order=scsi0
qm set ${TEMPLATE_ID} --serial0 socket --vga serial0

qm set ${TEMPLATE_ID} --agent enabled=1
qm set ${TEMPLATE_ID} --onboot 1

# echo "ssh-rsa " >authorized_keys
qm set ${TEMPLATE_ID} --ipconfig0 "ip6=auto,ip=dhcp"
# qm set ${TEMPLATE_ID} --sshkeys ~/authorized_keys

qm set ${TEMPLATE_ID} --ciuser ubuntu
qm set ${TEMPLATE_ID} --cipassword ubuntu@123
qm template ${TEMPLATE_ID}

echo "***********TEMPLATE ${TEMPLATE_NAME} successfully created!************" &&
    echo "Now create a clone of VM with ID ${TEMPLATE_ID} in the Webinterface.."

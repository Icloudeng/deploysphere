#!/bin/bash
export TEMPLATE_IMAGE="jammy-server-cloudimg-amd64"
export TEMPLATE_NAME="ubuntu-2204-cloudinit-template"
export TEMPLATE_ID=9000
export VM_DISK="local-lvm"
export VM_BRIDGE="vmbr0"
sudo apt update -y
sudo apt install libguestfs-tools -y
if [ ! -f "${TEMPLATE_IMAGE}.img" ]; then
    wget https://cloud-images.ubuntu.com/jammy/current/${TEMPLATE_IMAGE}.img
fi
sudo virt-customize -a ${TEMPLATE_IMAGE}.img --install qemu-guest-agent,ncat,mc,net-tools,bash-completion --run-command 'systemctl enable qemu-guest-agent.service'
cp ${TEMPLATE_IMAGE}.img ${TEMPLATE_IMAGE}.qcow2

qemu-img resize ${TEMPLATE_IMAGE}.qcow2 30G
qm destroy ${TEMPLATE_ID}
qm create ${TEMPLATE_ID} --name ${TEMPLATE_NAME} --memory 2048 --cores 2 --net0 virtio,bridge=${VM_BRIDGE}
qm importdisk ${TEMPLATE_ID} ${TEMPLATE_IMAGE}.qcow2 ${VM_DISK}
qm set ${TEMPLATE_ID} --scsihw virtio-scsi-single --scsi0 ${VM_DISK}:vm-${TEMPLATE_ID}-disk-0,cache=none,ssd=1,discard=on
qm set ${TEMPLATE_ID} --boot c --bootdisk scsi0
qm set ${TEMPLATE_ID} --ide2 ${VM_DISK}:cloudinit
qm set ${TEMPLATE_ID} --serial0 socket --vga serial0
qm set ${TEMPLATE_ID} --agent enabled=1
qm set ${TEMPLATE_ID} --onboot 1
# echo "ssh-rsa " >authorized_keys
qm set ${TEMPLATE_ID} --sshkeys ~/authorized_keys
qm set ${TEMPLATE_ID} --ipconfig0 "ip6=auto,ip=dhcp"
qm set ${TEMPLATE_ID} --ciuser ubuntu
qm set ${TEMPLATE_ID} --cipassword ubuntu@123
qm template ${TEMPLATE_ID}

echo "***********TEMPLATE ${TEMPLATE_NAME} successfully created!************" &&
    echo "Now create a clone of VM with ID ${TEMPLATE_ID} in the Webinterface.."

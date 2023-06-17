#!/bin/bash

# Default values
ansible_user=""
platform=""
vm_ip=""
metadata=""

# Parse named arguments
while [[ $# -gt 0 ]]; do
    key="$1"

    case $key in
    --ansible-user)
        ansible_user="$2"
        shift
        shift
        ;;
    --platform)
        platform="$2"
        shift
        shift
        ;;
    --vmip)
        vm_ip="$2"
        shift
        shift
        ;;
    --metadata)
        metadata="$2"
        shift
        shift
        ;;
    *)
        echo "Invalid argument: $1"
        exit 1
        ;;
    esac
done

# # Check if the arguments was provided
if [ -z "$ansible_user" ] || [ -z "$platform" ] || [ -z "$vm_ip" ]; then
    echo "Platform, vmip and ansible-user arguments are required. Usage: $0 --platform <platform> --vmip <vmip> --ansible-user <ansible-user>"
    exit 1
fi

# Check if Python or Python 3 is installed
if command -v python &>/dev/null; then
    echo "Python is installed."
    python_version=$(python --version 2>&1)
    echo "Python version: $python_version"
    python_command="python"
elif command -v python3 &>/dev/null; then
    echo "Python 3 is installed."
    python_version=$(python3 --version 2>&1)
    echo "Python 3 version: $python_version"
    python_command="python3"
else
    echo "Python or Python 3 is not installed. Please install Python and try again."
    exit 1
fi

# Check if the Ansible playbook script exists
playbook_path="./scripts/platforms/$platform/playbook.yaml"
if [ -f "$playbook_path" ]; then
    echo "Ansible playbook script found: $playbook_path"
else
    echo "Ansible playbook script not found for platform: $platform"
    exit 1
fi

# Create Python virtual environment
$python_command -m venv .venv
source .venv/bin/activate

# Set PLATFORM_INSTALLER_METADATA environment variable
export PLATFORM_INSTALLER_METADATA="$metadata"

# Upgrade pip
$python_command -m pip install --upgrade pip

# Check if Ansible is already installed
if $python_command -c "import ansible" &>/dev/null; then
    echo "Ansible is already installed."
else
    # Install Ansible in the virtual environment
    echo "Ansible was not found, start installing..."
    pip install ansible jmespath
fi

export ANSIBLE_HOST_KEY_CHECKING=False

# Run Ansible playbook
if [ -f "./private-key.pem" ]; then
    chmod 600 ./private-key.pem
    ansible-playbook -u "$ansible_user" -i "$vm_ip," --private-key "./private-key.pem" "$playbook_path"
else
    echo ansible-playbook -u "$ansible_user" -i "'$vm_ip,'" "$playbook_path"
    ansible-playbook -u "$ansible_user" -i "'$vm_ip,'" "$playbook_path"
fi

# Deactivate the virtual environment
deactivate

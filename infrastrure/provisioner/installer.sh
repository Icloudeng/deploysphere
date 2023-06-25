#!/bin/bash

MY_DIR=$(dirname $0)

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

# Check if the Ansible playbook script exists
playbook_path="./scripts/platforms/$platform/playbook.yaml"

if [ -f "$playbook_path" ]; then
    echo "Ansible playbook script found: $playbook_path"
else
    echo "Ansible playbook script not found for platform: $platform"
    exit 1
fi

# Include python command and activate python venv
source $MY_DIR/bash/ansible_init.sh

# Get the last total ansible logs file line number
logs_lines=$(wc -l <$ansible_log_file | tr -d '[:space:]')

# Read variables from /root/.env variable and pass them to extra variable
getenv="$python_command lib/getenv.py"

################ Ansible extra-vars ################
ansible_extra_vars="platform_metadata=$metadata"

# Run Ansible playbook
if [ -f "./private-key.pem" ]; then
    chmod 600 ./private-key.pem
    ansible-playbook -u "$ansible_user" -i "$vm_ip," --private-key "./private-key.pem" "$playbook_path" --extra-vars "$ansible_extra_vars"
    # Capture the exit code of the Ansible playbook command
    playbook_result=$?
else
    ansible-playbook -u "$ansible_user" -i "'$vm_ip,'" "$playbook_path" --extra-vars "$ansible_extra_vars"
    # Capture the exit code of the Ansible playbook command
    playbook_result=$?
fi

# Get the ansible logs content from last run and pipe it to base64
ansible_logs=$(tail -n +$logs_lines $ansible_log_file | base64)

# Send notification accordingly
if [ $playbook_result -eq 0 ]; then
    echo "Playbook succeeded!"
    ran_status="succeeded"
else
    echo "Playbook failed!"
    ran_status="failed"
fi

# Execute python notifier script
$python_command lib/notifier.py --logs "$ansible_logs" --status "$ran_status" --platform "$platform" --ip "$vm_ip"

# Execute Nginx Proxy Manager (Domain mapping)
$python_command lib/nginx_pm.py --action "create" --metadata "$metadata" --platform "$platform" --status "$ran_status" --ip "$vm_ip"

# Deactivate the virtual environment
deactivate

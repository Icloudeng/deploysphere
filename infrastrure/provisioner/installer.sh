#!/bin/bash

MY_DIR=$(dirname $0)
PROXY_MANAGER="nginx" # nginx | treafik

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
source $MY_DIR/bash/init.sh
source $MY_DIR/bash/functions.sh

# Generate static token based on platform name
static_secret=$(get_platform_static_secret)

# Get the last total ansible logs file line number
logs_lines=$(wc -l <$ansible_log_file | tr -d '[:space:]')

################ Ansible extra-vars ################
ansible_extra_vars="platform_metadata=$metadata platform_name=$platform"
ansible_extra_vars+=" random_secret=$random_secret admin_email=$admin_email" # Must start with empty space
ansible_extra_vars+=" static_secret=$static_secret"                          # Must start with empty space

# Notification Installer details
installer_details="Platform: $platform\nMachine IP: $vm_ip\n\n"
installer_details+="Static Secret=$static_secret\n"

# Notify before playbook
message_info=$(echo "Provisioning Started..." | base64)
$python_command lib/notifier.py --logs "$message_info" --status "info" --details "$installer_details" --metadata "$metadata"

# Run Ansible playbook (Function) -> (ran_status: succeeded | failed)
execute_ansible_playbook

# Create domain mapping with (nginx|treafik)
if [ "$ran_status" == "succeeded" ]; then

    if [ "$PROXY_MANAGER" == "nginx" ]; then
        # Execute Nginx Proxy Manager (Domain mapping)
        $python_command lib/nginx_pm.py --action "create" --metadata "$metadata" --platform "$platform" --ip "$vm_ip"
    elif [ "$PROXY_MANAGER" == "treafik" ]; then

        echo "Proxy manager with treafik"
    fi

fi

# Get the ansible logs content from last run and pipe it to base64
ansible_logs=$(tail -n +$logs_lines $ansible_log_file)
ansible_logs_4096=$(get_last_n_chars "$ansible_logs" 4096 | base64)

# Read and extract variables exposed from ansible logs
exposed_variables=$($extract_vars "$ansible_logs")

# Execute python notifier script
installer_details+="Random Secret=$random_secret\n\n$exposed_variables\n"
$python_command lib/notifier.py --logs "$ansible_logs_4096" --status "$ran_status" --details "$installer_details" --metadata "$metadata"

# if [ "$ran_status" == "succeeded" ]; then
#     # Include LDAP Script
#     source $MY_DIR/bash/ldap_executor.sh
# fi

# Deactivate the virtual environment
deactivate

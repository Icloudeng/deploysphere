#!/bin/bash

MY_DIR=$(dirname $0)
PROXY_MANAGER="nginx" # nginx | treafik

# Default values
ansible_user=""
platform=""
vm_ip=""
job_id=""
metadata=""
reference=""

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
    --job-id)
        job_id="$2"
        shift
        shift
        ;;
    --reference)
        reference="$2"
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
ansible_extra_vars+=" random_secret=$random_secret global_admin_email=$admin_email" # Must start with empty space
ansible_extra_vars+=" static_secret=$static_secret"                          # Must start with empty space
ansible_extra_vars+=" vm_ip=$vm_ip"                                          # Must start with empty space

# Notification Installer details
# @@Function@@ -> (var: installer_details)
fill_installer_details_installer

# Notify before playbook
message_info=$(echo "Provisioning Started..." | base64)
$python_command lib/notifier.py --logs "$message_info" --status "info" --details "$installer_details" --metadata "$metadata"

# @@Function@@ (Run Ansible playbook) -> (ran_status: succeeded | failed, channel_publisher)
execute_ansible_playbook

# @@Function@@  Create domain mapping with (nginx|treafik)
create_domain_mapping

# Get the ansible logs content from last run and pipe it to base64
ansible_logs=$(tail -n +$logs_lines $ansible_log_file)
ansible_logs_8191=$(get_last_n_chars "$ansible_logs" 8191)
ansible_logs_4096=$(get_last_n_chars "$ansible_logs" 4096 | base64)

# @@Function@@
publish_redis_playbook_details

# Read and extract variables exposed from ansible logs
exposed_variables=$($extract_vars --text "$ansible_logs_8191")
# Execute python notifier script
installer_details+="$exposed_variables\n"

$python_command lib/notifier.py --logs "$ansible_logs_4096" --status "$ran_status" --details "$installer_details" --metadata "$metadata"

# Deactivate the virtual environment
deactivate

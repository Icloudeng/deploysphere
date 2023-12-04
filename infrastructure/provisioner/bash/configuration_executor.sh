# Check if the Ansible playbook script exists
playbook_path="./scripts/platforms/$platform/playbook-configuration.yaml"

if [ -f "$playbook_path" ]; then
    echo "Ansible Configuration playbook script found: $playbook_path"
else
    echo "Ansible Configuration playbook script not found for platform: $platform"
    exit 1
fi

# Test configuration object
configuration_data=$($get_decoded_metadata --key "configuration" --metadata "$metadata")
if jq -e . >/dev/null 2>&1 <<<"$configuration_data"; then
    echo "Parsed Configuration JSON successfully"
else
    echo "Failed to parse Configuration JSON"
    exit 1
fi

# Overwrite global varaible
# Generate static token based on platform name
static_secret=$(get_platform_static_secret)

################ Ansible extra-vars ################
ansible_extra_vars="platform_metadata=$metadata platform_name=$platform"
ansible_extra_vars+=" random_secret=$random_secret global_admin_email=$admin_email" # Must start with empty space
ansible_extra_vars+=" static_secret=$static_secret"                          # Must start with empty space
ansible_extra_vars+=" vm_ip=$vm_ip"                                          # Must start with empty space

# Get the last total ansible logs file line number (Before ansible execution)
logs_lines=$(wc -l <$ansible_log_file | tr -d '[:space:]')

# Run Ansible playbook (Function) -> (ran_status: succeeded | failed)
execute_ansible_playbook

# Get the ansible logs content from last run and pipe it to base64 (After ansible execution)
ansible_logs=$(tail -n +$logs_lines $ansible_log_file)
ansible_logs_8191=$(get_last_n_chars "$ansible_logs" 8191)
ansible_logs_4096=$(get_last_n_chars "$ansible_logs" 4096 | base64)

# @@Function@@
publish_redis_playbook_details

# Read and extract variables exposed from ansible logs
exposed_variables=$($extract_vars --text "$ansible_logs_8191")

# @@Function@@ -> (var: installer_details)
fill_installer_details_configuration

# Execute python notifier script

$python_command lib/notifier.py --logs "$ansible_logs_4096" --status "$ran_status" --details "$installer_details" --metadata "$metadata"

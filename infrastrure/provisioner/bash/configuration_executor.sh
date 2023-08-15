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
ansible_extra_vars+=" random_secret=$random_secret admin_email=$admin_email" # Must start with empty space
ansible_extra_vars+=" static_secret=$static_secret"                          # Must start with empty space
ansible_extra_vars+=" vm_ip=$vm_ip"                                          # Must start with empty space

# Get the last total ansible logs file line number (Before ansible execution)
logs_lines=$(wc -l <$ansible_log_file | tr -d '[:space:]')

# Run Ansible playbook (Function) -> (ran_status: succeeded | failed)
execute_ansible_playbook

# Get the ansible logs content from last run and pipe it to base64 (After ansible execution)
ansible_logs=$(tail -n +$logs_lines $ansible_log_file)
ansible_logs_4096=$(get_last_n_chars "$ansible_logs" 4096 | base64)

# Publish Credentials
if [ "$ran_status" == "succeeded" ]; then
    # Read and extract credentials exposed from ansible logs
    exposed_credentials=$($extract_vars --text "$ansible_logs" --credentials "true")

    # Publish credentials if not empty
    if [ -n "$exposed_credentials" ]; then
        $redis_publisher --channel "$channel_publisher-credentials" --message "$(echo "$exposed_credentials" | base64)"
    fi

fi

# Publish playbook run status
$redis_publisher --channel "$channel_publisher-status" --message "$(echo "$ran_status" | base64)"

# Read and extract variables exposed from ansible logs
exposed_variables=$($extract_vars --text "$ansible_logs")

# Execute python notifier script
installer_details="EXECUTION TYPE: Configuration\n"
installer_details+="Reference: $reference\n\n"

installer_details+="Platform: $platform\n"
installer_details+="Machine User: $ansible_user\nMachine IP: $vm_ip\n\n"
installer_details+="$exposed_variables"

$python_command lib/notifier.py --logs "$ansible_logs_4096" --status "$ran_status" --details "$installer_details" --metadata "$metadata"

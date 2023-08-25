# Function to get the last n characters from a string
get_last_n_chars() {
    local input_string="$1"
    local n_chars="$2"

    # Get the length of the string
    local string_length=${#input_string}

    # Check if the string length is less than or equal to n_chars
    if ((string_length <= n_chars)); then
        echo "$input_string"
    else
        # Calculate the starting position for the last n_chars
        local start_position=$((string_length - n_chars))

        # Extract the last n_chars using parameter expansion
        echo "${input_string:$start_position}"
    fi
}

execute_ansible_playbook() {
    if [ -n "$reference" ]; then
        channel_publisher="$reference"
    else
        channel_publisher="$vm_ip"
    fi

    if [ -f "./private-key.pem" ]; then
        chmod 600 ./private-key.pem
        # Run the command in background
        ansible-playbook -u "$ansible_user" -i "$vm_ip," --private-key "./private-key.pem" "$playbook_path" --extra-vars "$ansible_extra_vars" --extra-vars "@scripts/variables.yaml" &
        # capture the process ID of the Ansible playbook command
        playbook_pid=$!

        # =====================

        # Run the logs exports
        $logs_publisher --channel "$channel_publisher" &
        # capture the process ID of the Logs exporter
        logs_publisher_pid=$!
    else
        # Run the command in background
        ansible-playbook -u "$ansible_user" -i "'$vm_ip,'" "$playbook_path" --extra-vars "$ansible_extra_vars" --extra-vars "@scripts/variables.yaml" &
        # capture the process ID of the Ansible playbook command
        playbook_pid=$!

        # =====================

        # Run the logs publisher
        $logs_publisher --channel "$channel_publisher" &
        # capture the process ID of the Logs exporter
        logs_publisher_pid=$!
    fi

    # Wait for both commands to finish and capture their exit codes
    wait $playbook_pid
    # Capture the exit code of the Ansible playbook command
    playbook_result=$?

    # Wait just two second to make sure all python exportation netword process has been publish
    sleep 5

    # Kill the logs publisher process
    kill -SIGINT $logs_publisher_pid >/dev/null 2>&1

    if [ $playbook_result -eq 0 ]; then
        echo "Playbook succeeded!"
        ran_status="succeeded"
    else
        echo "Playbook failed!"
        ran_status="failed"
    fi
}

get_platform_static_secret() {
    static_secret_name="$platform-$(date +%Y-%m)"
    static_secret=$($python_command -c "import hashlib; print(hashlib.sha256('$static_secret_name'.encode()).hexdigest()[:32])")
    echo $static_secret
}

create_domain_mapping() {
    if [ "$ran_status" == "succeeded" ]; then

        if [ "$PROXY_MANAGER" == "nginx" ]; then
            # Execute Nginx Proxy Manager (Domain mapping)
            $python_command lib/nginx_pm.py --action "create" --metadata "$metadata" --platform "$platform" --ip "$vm_ip"
        elif [ "$PROXY_MANAGER" == "treafik" ]; then

            echo "Proxy manager with treafik"
        fi

    fi

}

publish_redis_playbook_details() {
    # Publish Credentials
    if [ "$ran_status" == "succeeded" ]; then
        # Read and extract credentials exposed from ansible logs
        exposed_credentials=$($extract_vars --text "$ansible_logs" --credentials "true")

        # Publish credentials if not empty
        if [ -n "$exposed_credentials" ]; then
            $redis_publisher --channel "$channel_publisher-credentials" --message "$exposed_credentials"
        fi

    fi

    # Publish playbook run status
    $redis_publisher --channel "$channel_publisher-status" --message "$ran_status"
}

fill_installer_details_installer() {
    installer_details=""

    if [ -n "$reference" ]; then
        installer_details+="Reference: $reference\n\n"
    fi

    if [ -n "$job_id" ]; then
        installer_details+="Job ID: $job_id\n"
    fi

    installer_details+="Platform: $platform\n"
    installer_details+="Machine User: $ansible_user\n"
    installer_details+="Machine IP: $vm_ip\n\n"

    installer_details+="Static Secret: $static_secret\n"
    installer_details+="Random Secret: $random_secret\n\n"
}

fill_installer_details_configuration() {
    installer_details="EXECUTION TYPE: Configuration\n"

    if [ -n "$reference" ]; then
        installer_details+="Reference: $reference\n\n"
    fi

    if [ -n "$job_id" ]; then
        installer_details+="Job ID: $job_id\n"
    fi

    installer_details+="Platform: $platform\n"
    installer_details+="Machine User: $ansible_user\nMachine IP: $vm_ip\n\n"
    installer_details+="$exposed_variables"
}

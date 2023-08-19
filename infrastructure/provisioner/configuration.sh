#!/bin/bash

MY_DIR=$(dirname $0)

# Default values
ansible_user=""
platform=""
vm_ip=""
metadata=""
reference=""
job_id=""

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

# Include python command and activate python venv
source $MY_DIR/bash/init.sh
source $MY_DIR/bash/functions.sh

# Include configuration Script
source $MY_DIR/bash/configuration_executor.sh

# Deactivate the virtual environment
deactivate

#!/bin/bash

MY_DIR=$(dirname $0)
# Default values
metadata=""
action=""
platform=""
ip=""

# Parse named arguments
while [[ $# -gt 0 ]]; do
    key="$1"

    case $key in
    --action)
        action="$2"
        shift
        shift
        ;;
    --metadata)
        metadata="$2"
        shift
        shift
        ;;
    --platform)
        platform="$2"
        shift
        shift
        ;;
    --ip)
        ip="$2"
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
if [ -z "$metadata" ] || [ -z "$action" ]; then
    echo "metadata, action argument required. Usage: $0 --action <delete|create> --metadata <metadata>"
    exit 1
fi

# Include python command and activate python venv
source $MY_DIR/bash/ansible_init.sh

# Execute Nginx Proxy Manager (Domain mapping)
$python_command lib/nginx_pm.py --action "$action" --metadata "$metadata" --platform "$platform" --ip "$ip"

# Deactivate the virtual environment
deactivate

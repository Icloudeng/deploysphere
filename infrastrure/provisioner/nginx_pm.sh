#!/bin/bash

MY_DIR=$(dirname $0)
# Default values
metadata=""

# Parse named arguments
while [[ $# -gt 0 ]]; do
    key="$1"

    case $key in
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
if [ -z "$metadata" ]; then
    echo "metadata argument required. Usage: $0 --metadata <metadata>"
    exit 1
fi

# Include python command and activate python venv
source $MY_DIR/bash/ansible_init.sh

# Execute Nginx Proxy Manager (Domain mapping)
$python_command lib/nginx_pm.py --action "delete" --metadata "$metadata"

# Deactivate the virtual environment
deactivate

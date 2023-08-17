#!/bin/bash

MY_DIR=$(dirname $0)
# Default values
logs=""
status=""
details=""
metadata=""

# Parse named arguments
while [[ $# -gt 0 ]]; do
    key="$1"

    case $key in
    --logs)
        logs="$2"
        shift
        shift
        ;;
    --status)
        status="$2"
        shift
        shift
        ;;
    --details)
        details="$2"
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
if [ -z "$logs" ] || [ -z "$status" ]; then
    echo "logs, status argument required. Usage: $0 --status <status> --logs <logs>"
    exit 1
fi

# Include python command and activate python venv
source $MY_DIR/bash/init.sh

$python_command lib/notifier.py --logs "$logs" --status "$status" --details "$details" --metadata "$metadata"

# Deactivate the virtual environment
deactivate

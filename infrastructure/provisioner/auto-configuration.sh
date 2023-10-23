#!/bin/bash

MY_DIR=$(dirname $0)

# Default values
type=""
platform=""
reference=""
config_reference=""

# Parse named arguments
while [[ $# -gt 0 ]]; do
    key="$1"

    case $key in
    --type)
        type="$2"
        shift
        shift
        ;;
    --platform)
        platform="$2"
        shift
        shift
        ;;
    --config-reference)
        config_reference="$2"
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
if [ -z "$type" ] || [ -z "$platform" ] || [ -z "$config_reference" ] || [ -z "$reference" ]; then
    echo "type, platform, reference, config-reference arguments are required. Usage: $0 --type <type> --platform <platform> ..."
    exit 1
fi

# Include python command and activate python venv
source $MY_DIR/bash/init.sh
source $MY_DIR/bash/functions.sh

# RUN Python script
$python_command scripts/configurations/$platform/$type.py --reference "$reference" --config-reference "$config_reference" --type "$type" --platform "$platform"
# Capture the exit code
execution_code=$?

# Deactivate the virtual environment
deactivate

exit $execution_code

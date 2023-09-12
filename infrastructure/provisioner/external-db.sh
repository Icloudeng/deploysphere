#!/bin/bash

MY_DIR=$(dirname $0)

# Default values
db_type=""
db_name=""

# Parse named arguments
while [[ $# -gt 0 ]]; do
    key="$1"

    case $key in
    --db-type)
        db_type="$2"
        shift
        shift
        ;;
    --db-name)
        db_name="$2"
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
if [ -z "$db_type" ] || [ -z "$db_name" ]; then
    echo "db-name, db-type arguments are required. Usage: $0 --db-type <db-type> --db-name <db-name>"
    exit 1
fi

# Include python command and activate python venv
source $MY_DIR/bash/init.sh

$external_db --db-type "$db_type" --db-name "$db_name"
# Capture the exit code
execution_code=$?

# Deactivate the virtual environment
deactivate

exit $execution_code

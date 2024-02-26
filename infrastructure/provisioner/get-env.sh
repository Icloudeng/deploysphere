#!/bin/bash

MY_DIR=$(dirname $0)

# The first argument is stored in $1
first_argument=$1

# Check if an argument was provided
if [ -z "$first_argument" ]; then
    echo "Please provide an argument."
    exit 1
fi

# Include python command and activate python venv
source $MY_DIR/bash/init.sh

$getenv "$first_argument"
# Capture the exit code
execution_code=$?

# Deactivate the virtual environment
deactivate

exit $execution_code

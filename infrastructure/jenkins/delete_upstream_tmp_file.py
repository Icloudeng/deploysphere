#!/usr/bin/env python3

import re
import os
import tempfile

BUILD_TAG = os.environ.get("BUILD_TAG", "")


def format_to_filename_standard(input_string):
    # Remove leading and trailing whitespace
    input_string = input_string.strip()

    # Replace spaces and other non-alphanumeric characters with hyphens
    formatted_string = re.sub(r"[^\w\s-]", "", input_string)
    formatted_string = re.sub(r"[\s]+", "-", formatted_string)

    return formatted_string


def delete_file():
    # Get the temporary directory path
    temp_dir = tempfile.gettempdir()
    # Specify the filename
    file_path = os.path.join(temp_dir, format_to_filename_standard(BUILD_TAG))

    try:
        if os.path.isfile(file_path):
            os.remove(file_path)
            print(f"File '{file_path}' has been deleted.")
        else:
            print(f"File '{file_path}' does not exist.")
    except Exception as e:
        print(f"An error occurred while deleting the file: {str(e)}")


if __name__ == "__main__":
    if BUILD_TAG:
        delete_file()

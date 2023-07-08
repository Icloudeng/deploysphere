import re
import sys


def extract_content_all(input_string: str):
    # Regular expression pattern
    pattern = r"%%%(.*?)%%%"

    # Extract all content using regex
    content_array = re.findall(pattern, input_string)

    # Return extracted content
    return content_array


if __name__ == '__main__':
    extracted_content = extract_content_all(sys.argv[1].strip())
    print("\n".join(extracted_content), end="")

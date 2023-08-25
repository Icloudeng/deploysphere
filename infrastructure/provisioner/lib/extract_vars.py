import argparse
import re


def extract_content_all(input_string: str, pattern: str):
    # Extract all content using regex
    content_array = re.findall(pattern, input_string)

    # Return extracted content
    return content_array


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument("--text", required=True)
    parser.add_argument(
        "--credentials",
        required=False,
        default=False,
        type=bool
    )
    args = parser.parse_args()

    # Regular expression pattern
    pattern = r"%\$%(.*?)%\$%" if args.credentials else r"%%%(.*?)%%%"

    extracted_content = extract_content_all(
        input_string=args.text,
        pattern=pattern
    )

    if args.credentials:
        content = extracted_content.replace("\\", "")
        content = '[%s]' % (','.join(extracted_content))
    else:
        content = "\n".join(extracted_content)

    print(content, end="")

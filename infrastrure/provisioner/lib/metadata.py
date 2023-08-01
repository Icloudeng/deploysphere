

import argparse
import base64
import json
from utilities.logging import logging, bingLoggingConfig


def decode_metadata(metadata: str):
    data = {}
    try:
        decoded_bytes = base64.b64decode(metadata)
        data = json.loads(decoded_bytes.decode("utf-8"))
    except Exception as e:
        logging.warning(e)

    return data


def get_metadata(key, metadata: str):
    decoded_matadata = decode_metadata(metadata)

    if key and len(key) > 0:
        decoded_matadata = decoded_matadata.get(key, "")

    return json.dumps(decoded_matadata)


if __name__ == '__main__':
    bingLoggingConfig(prefix="Metadata / ")

    parser = argparse.ArgumentParser()
    parser.add_argument("--metadata", required=True, default="")
    parser.add_argument("--key", required=False, default="")
    args = parser.parse_args()

    try:
        print(get_metadata(key=args.key, metadata=args.metadata), end="")
    except:
        print("", end="")

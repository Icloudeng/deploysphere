import argparse
from utilities.logging import logging, bingLoggingConfig


if __name__ == '__main__':
    bingLoggingConfig(prefix="Logger /")

    parser = argparse.ArgumentParser()
    parser.add_argument("--text", required=True, default="")
    args = parser.parse_args()

    logging.debug(args.text)

import base64
import sys
import argparse
import redis
from utilities.dotenv import config
from utilities.logging import logging, bingLoggingConfig


def main(channel: str, message: str):
    REDIS_URL: str = config.get("REDIS_URL")
    redis_url = REDIS_URL.split(":")
    redis_client = redis.Redis(host=redis_url[0], port=int(redis_url[1]), db=0)

    encoded_message = base64.b64encode(
        message.strip().encode('utf-8')
    ).decode('utf-8')

    redis_client.publish(channel=channel, message=encoded_message)

    redis_client.close()

    sys.exit(0)


if __name__ == '__main__':
    bingLoggingConfig(prefix="Redis Exporter / ")

    parser = argparse.ArgumentParser()
    parser.add_argument("--channel", required=True)
    parser.add_argument("--message", required=True)
    args = parser.parse_args()

    logging.info(args)

    main(channel=args.channel, message=args.message)

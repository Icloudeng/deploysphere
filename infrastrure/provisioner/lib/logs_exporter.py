import signal
import sys
import subprocess
import argparse
import redis
from utilities.dotenv import config
from utilities.logging import logging, bingLoggingConfig

log_file = "./logs/ansible_log.txt"


def main(channel: str):
    CHANNEL = f"{channel}-logs"

    REDIS_URL: str = config.get("REDIS_URL")
    redis_url = REDIS_URL.split(":")
    redis_client = redis.Redis(host=redis_url[0], port=int(redis_url[1]), db=0)

    tail_process = subprocess.Popen(
        ['tail', '-n', '-0', '-f', log_file],
        stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True
    )

    def signal_handler(sig, frame):
        print("Script is exiting")
        tail_process.terminate()
        redis_client.close()
        sys.exit(0)

    # Register signal handlers
    signal.signal(signal.SIGINT, signal_handler)  # Ctrl+C
    signal.signal(signal.SIGTERM, signal_handler)  # Termination signal

    try:
        while True:
            # Read new lines from the tail process
            new_line = tail_process.stdout.readline()
            if new_line:
                redis_client.publish(CHANNEL, new_line.strip())
    except KeyboardInterrupt:
        # Stop the tail process if the user presses Ctrl+C
        tail_process.terminate()


if __name__ == '__main__':
    bingLoggingConfig(prefix="Ansible Logs Exporter / ")

    parser = argparse.ArgumentParser()
    parser.add_argument("--channel")
    args = parser.parse_args()

    logging.info(args)

    main(channel=args.channel)

import logging


def bingLoggingConfig(prefix: str = ''):
    root_logger = logging.getLogger()
    root_logger.setLevel(logging.DEBUG)  # or whatever
    handler = logging.FileHandler(
        'logs/python-lib.log', 'w', 'utf-8'
    )  # or whatever
    handler.setFormatter(
        logging.Formatter(
            f'{prefix}%(name)s - %(levelname)s - %(message)s'
        )
    )  # or whatever
    root_logger.addHandler(handler)

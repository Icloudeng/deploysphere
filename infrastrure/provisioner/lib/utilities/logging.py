import logging


def bingLoggingConfig(prefix: str = ''):
    logging.basicConfig(
        filename='python-lib.log',
        encoding='utf-8',
        level=logging.DEBUG,
        format=f'{prefix}%(name)s - %(levelname)s - %(message)s'
    )

import sys
from utilities.dotenv import config


if __name__ == '__main__':
    print(config.get(sys.argv[1].strip(), ""), end="")

import base64
import json
import argparse
import re
from utilities.dotenv import config
from utilities.logging import logging, bingLoggingConfig
from sqlalchemy import create_engine
from sqlalchemy_utils import database_exists, create_database


PREFIX = "EXTERNAL_"

databases_fields = {
    "postgres": ["HOST", "PORT", "USER", "PASSWORD", "TIMEZONE", "SSLMODE"],
    "mysql": ["HOST", "PORT", "USER", "PASSWORD"],
}

protocols = {"postgres": "postgresql", "mysql": "mysql+pymysql"}


def format_database_name(input_str: str):
    # Remove any characters that are not letters, numbers, or underscores
    sanitized_str = re.sub(r"[^a-zA-Z0-9_]", "", input_str)

    # Ensure the name starts with a letter or an underscore
    if not sanitized_str[0].isalpha() and sanitized_str[0] != "_":
        sanitized_str = "_" + sanitized_str

    # Ensure the name is no longer than 63 characters
    if len(sanitized_str) > 63:
        sanitized_str = sanitized_str[:63]

    return sanitized_str


def db_connection(db_type: str, db_name: str, fields):
    protocol = protocols.get(db_type)

    engine = create_engine(
        f"{protocol}://{fields['USER']}:{fields['PASSWORD']}@{fields['HOST']}:{fields['PORT']}/{db_name}"
    )

    if not database_exists(engine.url):
        create_database(engine.url)
        print(f"Database '{db_name}' created successfully.")


def main(db_type: str, db_name: str):
    db_fields = databases_fields.get(db_type, None)
    if db_fields == None:
        raise Exception(
            f"\n\nThe requested database configuration is not found: {db_type}\n\n"
        )

    fields = {"NAME": db_name}
    for field in db_fields:
        field_idx = PREFIX + f"{db_type.upper()}_" + field
        value = config.get(field_idx, None)

        if not value:
            raise Exception(
                f"\n\nEmpty Field Database Configuration: database type ({db_type}), field ({field_idx}) \n\n"
            )

        fields[field] = value

    db_connection(db_type, db_name, fields)

    return fields


if __name__ == "__main__":
    bingLoggingConfig(prefix="External Database Access / ")

    parser = argparse.ArgumentParser()
    parser.add_argument("--db-type", required=True)
    parser.add_argument("--db-name", required=True)
    args = parser.parse_args()

    logging.info(args)

    # Format database name
    db_name = format_database_name(args.db_name)

    data = main(db_name=db_name, db_type=args.db_type)

    json_bytes = json.dumps(data).encode("utf-8")

    base64_encoded = base64.b64encode(json_bytes).decode("utf-8")

    print(f"%%%${base64_encoded}%%%")

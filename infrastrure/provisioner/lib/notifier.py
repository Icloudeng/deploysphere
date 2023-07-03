import json
import base64
import asyncio
import argparse
import telegram
from utilities.dotenv import config
from utilities.logging import logging, bingLoggingConfig

DOMAIN_KEY = "domain"


def get_message_content(base64_content: str):
    # Decode the Base64-encoded string
    decoded_bytes = base64.b64decode(base64_content)
    # Convert bytes to string
    return decoded_bytes.decode("utf-8")


def create_bot():
    bot_token = config.get('TELEGRAM_BOT_TOKEN', None)
    chat_id = config.get('TELEGRAM_CHAT_ID', None)
    if not bot_token or not chat_id:
        return None, None

    # Create a Bot instance
    bot = telegram.Bot(token=bot_token)

    return bot, chat_id


def get_decoded_domain(metadata: str):
    domain = None
    try:
        decoded_bytes = base64.b64decode(metadata)
        data = json.loads(decoded_bytes.decode("utf-8"))
        domain = data.get(DOMAIN_KEY, None)
    except Exception as e:
        logging.warn(e)

    return domain


async def send_notification(encode_logs: str, status: str, installer_details: str, metadata: str):
    bot, chat_id = create_bot()
    if bot == None:
        logging.warn("Invalid BOT Configuration!")
        return
    # Decode message and send
    try:
        status_emoji = {
            "info": "🔵",
            "succeeded": "✅",
            "failed": "❌",
        }
        emoji = status_emoji.get(status, "🔵")

        decoded_logs = get_message_content(encode_logs).replace("--", "")

        sumzy = decoded_logs.find('========================================')
        sumzy = sumzy if sumzy > -1 else 0
        decoded_logs = decoded_logs[sumzy:]
        domain = get_decoded_domain(metadata)

        domain_text = f"\n\nDomain: {domain}\n\n" if domain else ""

        installer_details = installer_details.replace("\\n", "\n")

        content = f"##########################\n{decoded_logs[-3000:]}\n########################"
        text = f"\n{emoji} {status.title()}{domain_text}\n\n{installer_details}\n\n{content}"
        await bot.send_message(
            chat_id=chat_id,
            text=text
        )
    except Exception as error:
        logging.error("Failed to send notication", error)


if __name__ == '__main__':
    bingLoggingConfig(prefix="Notifier / ")

    parser = argparse.ArgumentParser()
    parser.add_argument("--logs", required=True)
    parser.add_argument("--status", required=True)
    parser.add_argument("--details", required=False, default="")
    parser.add_argument("--metadata", required=False, default="")
    args = parser.parse_args()

    logging.info(args)

    # Send notification
    asyncio.run(send_notification(
        encode_logs=args.logs,
        status=args.status,
        installer_details=args.details,
        metadata=args.metadata,
    ))
import sys
import base64
import asyncio
import argparse
import telegram
from .utilities.dotenv import config


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


async def send_notification(encode_logs: str, status: str, platform: str, ip: str):
    bot, chat_id = create_bot()
    if bot == None:
        print("Invalid BOT Configuration!")
        return
    # Decode message and send
    try:
        emoji = '❌' if status == "failed" else '✅'
        decoded_logs = get_message_content(encode_logs).replace("--", "")
        content = f"##########################\n{decoded_logs[-3500:]}\n########################"
        details = f"Platform: {platform}\nMachine IP: {ip}"
        text = f"{emoji} {status.title()}\n\n{details}\n\n{content}"
        await bot.send_message(
            chat_id=chat_id,
            text=text
        )
    except:
        print("Failed to send notication")


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument("--logs", required=True)
    parser.add_argument("--status", required=True)
    parser.add_argument("--platform", required=True)
    parser.add_argument("--ip", required=True)
    args = parser.parse_args()

    # Send notification
    asyncio.run(send_notification(
        encode_logs=args.logs,
        status=args.status,
        platform=args.platform,
        ip=args.ip
    ))

import os
import sys
import base64
import asyncio
import telegram
from dotenv import dotenv_values

# Get the current working directory
current_dir = os.getcwd()

# Get the parent folder (1 level up)
parent_dir = os.path.dirname(current_dir)

# Get the grandparent folder (2 levels up)
grandparent_dir = os.path.dirname(parent_dir)

# Specify the path to the .env file
dotenv_path = os.path.join(grandparent_dir, '.env')


config = dotenv_values(dotenv_path)


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


async def send_notification(encode_message: str, m_type: str, platform: str, ip: str):
    bot, chat_id = create_bot()
    if bot == None:
        print("Invalid BOT Configuration!")
        return
    # Decode message and send
    decoded_message = None
    try:
        emoji = '❌' if m_type == "failed" else '✅'
        decoded_message = get_message_content(encode_message)
        content = f"```bash\n##########################\n{decoded_message}\n########################```"
        details = f"Platform: {platform}\nMachine IP: {ip}"
        text = f"{emoji} *{m_type.title()}*\n\n{details}\n\n{content}"
        await bot.send_message(
            chat_id=chat_id,
            text=text,
            parse_mode=telegram.constants.ParseMode.MARKDOWN_V2
        )
    except:
        print("Failed to send notication")


if __name__ == '__main__':
    # Send notification
    asyncio.run(send_notification(
        encode_message=sys.argv[1],
        m_type=sys.argv[2],
        platform=sys.argv[3],
        ip=sys.argv[4]
    ))

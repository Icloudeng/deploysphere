import os
import sys
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

# Check if Python or Python 3 is installed
if command -v python &>/dev/null; then
    echo "Python is installed."
    python_version=$(python --version 2>&1)
    echo "Python version: $python_version"
    python_command="python"
elif command -v python3 &>/dev/null; then
    echo "Python 3 is installed."
    python_version=$(python3 --version 2>&1)
    echo "Python 3 version: $python_version"
    python_command="python3"
else
    echo "Python or Python 3 is not installed. Please install Python and try again."
    exit 1
fi

# Check if the OS is Ubuntu
if [[ "$(uname)" == "Linux" ]]; then
    # Check if pip is installed
    sudo apt-get update
    sudo apt-get -y install python3-venv python3-pip
fi

# Create Python virtual environment
$python_command -m venv .venv

if [ ! -f ".venv/bin/activate" ]; then
    echo "Venv script not found"
    exit 1
fi

source .venv/bin/activate

# Upgrade pip
$python_command -m pip install --upgrade pip

# Check if Ansible is already installed
if $python_command -c "import ansible, jmespath, telegram, dotenv, requests" &>/dev/null; then
    echo "Ansible is already installed."
else
    # Install Ansible in the virtual environment
    echo "Ansible was not found, start installing..."
    pip install ansible jmespath python-telegram-bot python-dotenv requests
fi

export ANSIBLE_HOST_KEY_CHECKING="False"
export ANSIBLE_CONFIG="$(pwd)/ansible.cfg"

ansible_log_file="ansible_log.txt"

# Create ansible log file if not exists
if [[ ! -f $ansible_log_file ]]; then
    touch "$ansible_log_file"
    echo "Created $ansible_log_file file."
fi

random_secret=$($python_command -c 'import secrets; print(secrets.token_hex(16))')

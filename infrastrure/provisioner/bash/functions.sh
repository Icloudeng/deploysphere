# Function to get the last n characters from a string
get_last_n_chars() {
    local input_string="$1"
    local n_chars="$2"

    # Get the length of the string
    local string_length=${#input_string}

    # Check if the string length is less than or equal to n_chars
    if ((string_length <= n_chars)); then
        echo "$input_string"
    else
        # Calculate the starting position for the last n_chars
        local start_position=$((string_length - n_chars))

        # Extract the last n_chars using parameter expansion
        echo "${input_string:$start_position}"
    fi
}

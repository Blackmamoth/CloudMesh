#!/bin/sh

# Usage information
usage() {
  echo "Usage: $0 <input_yaml_file> <output_file>"
  echo "Reads a YAML file and replaces placeholders with environment variable values"
  exit 1
}

# Check if we have required arguments
if [ $# -ne 2 ]; then
  usage
fi

input_file="$1"
output_file="$2"

# Check if input file exists
if [ ! -f "$input_file" ]; then
  echo "Error: Input file '$input_file' not found."
  exit 1
fi

# Create a temporary file
temp_file=$(mktemp)

# Copy input to temp file
cat "$input_file" >"$temp_file"

# Find all ${ENV_VAR} patterns in the input file
env_vars=$(grep -o '\${[A-Za-z0-9_]*}' "$input_file" | sort | uniq | sed 's/\${//;s/}//')

# Replace each environment variable
for var in $env_vars; do
  # Get the value of the environment variable
  value=$(eval echo \$$var)

  # If the environment variable is not set, warn the user
  if [ -z "$value" ]; then
    echo "Error: Environment variable $var is not set. Leaving placeholder unchanged."
    exit 1
  else
    # Replace the placeholder with the actual value
    # Using sed with different delimiters to avoid issues with slashes in values
    sed -i.bak "s|\${$var}|$value|g" "$temp_file"
  fi
done

# Move the temp file to the output location
mv "$temp_file" "$output_file"

echo "Successfully processed $input_file and saved to $output_file"
exit 0

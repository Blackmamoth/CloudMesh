#!/bin/bash

# Use awk to process the .env file and execute docker secret create for each line
awk -F= '
  # Skip empty lines and comments (lines starting with #)
  !/^$/ && !/^#/ {
    VAR = $1
    VAL = $2
    system("printf \"%s\" \"" VAL "\" | docker secret create " VAR " -")
  }
' .env.prod

#!/bin/bash

# To add this script to cron:
# 1. Make sure it's executable: chmod +x /path/to/api_checker.sh
# 2. Add it to cron: */5 7-23 * * * /path/to/api_checker.sh

# Get the directory where this script is located
PROJECT_DIR=$(dirname "$0")

# Change to that directory
cd $PROJECT_DIR

# Run the Python script using the virtual environment's Python
$PROJECT_DIR/venv/bin/python3 $PROJECT_DIR/api_checker.py

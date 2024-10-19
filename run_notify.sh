#!/bin/bash

# This script runs notify.py using the virtual environment's Python.
# To add this script to cron:
# 1. Make sure it's executable: chmod +x /path/to/run_notify.sh
# 2. Add it to cron: */5 7-23 * * * /path/to/run_notify.sh

# Get the directory where this script is located
PROJECT_DIR=$(dirname "$0")

# Change to that directory
cd $PROJECT_DIR

# Run the Python script using the virtual environment's Python
$PROJECT_DIR/venv/bin/python3 $PROJECT_DIR/notify.py >> $PROJECT_DIR/notify.log 2>&1

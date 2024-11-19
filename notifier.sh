#!/bin/bash

# Get the directory where this script is located
PROJECT_DIR=$(dirname "$0")

# Change to that directory
cd $PROJECT_DIR

# Run the Python script using the virtual environment's Python
$PROJECT_DIR/venv/bin/python3 $PROJECT_DIR/notifier.py

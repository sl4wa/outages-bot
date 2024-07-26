# Variables
PYTHON=python3
PIP=pip3
VENV=venv
SERVICE_NAME=poweron
SERVICE_FILE=$(SERVICE_NAME).service
SYSTEMD_SERVICE_FILE=/etc/systemd/system/$(SERVICE_NAME).service
USER=$(shell whoami)
BOT_DIR=$(shell pwd)
BOT_SCRIPT=main.py
GIT_REPO=https://github.com/sl4wa/poweron-bot.git
BRANCH=main

# Default target
.PHONY: all
all: help

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  setup          - Set up the Python virtual environment and install dependencies"
	@echo "  install        - Create the systemd service file and enable the service"
	@echo "  start          - Start the Telegram bot service"
	@echo "  stop           - Stop the Telegram bot service"
	@echo "  restart        - Restart the Telegram bot service"
	@echo "  status         - Show the status of the Telegram bot service"
	@echo "  logs           - Show the logs of the Telegram bot service"
	@echo "  uninstall      - Remove the systemd service file and disable the service"
	@echo "  update         - Pull the latest code from GitHub and restart the service"

# Clone the repository from GitHub
.PHONY: clone
clone:
	@if [ ! -d "$(BOT_DIR)" ]; then git clone -b $(BRANCH) $(GIT_REPO) $(BOT_DIR); fi

# Set up the Python virtual environment and install dependencies
.PHONY: setup
setup: clone
	cd $(BOT_DIR) && $(PYTHON) -m venv $(VENV)
	cd $(BOT_DIR) && $(VENV)/bin/$(PIP) install -r requirements.txt

# Create the systemd service file and enable the service
.PHONY: install
install: $(SYSTEMD_SERVICE_FILE)
	sudo systemctl daemon-reload
	sudo systemctl enable $(SERVICE_NAME)
	@echo "Service installed and enabled."

$(SYSTEMD_SERVICE_FILE): $(SERVICE_FILE)
	sudo mv $(SERVICE_FILE) $(SYSTEMD_SERVICE_FILE)

$(SERVICE_FILE): $(VENV)/bin/python
	@echo "Generating service file with virtual environment Python"
	@echo "[Unit]" > $(SERVICE_FILE)
	@echo "Description=Telegram Bot" >> $(SERVICE_FILE)
	@echo "After=network.target" >> $(SERVICE_FILE)
	@echo "" >> $(SERVICE_FILE)
	@echo "[Service]" >> $(SERVICE_FILE)
	@echo "User=$(USER)" >> $(SERVICE_FILE)
	@echo "WorkingDirectory=$(BOT_DIR)" >> $(SERVICE_FILE)
	@echo "ExecStart=$(BOT_DIR)/$(VENV)/bin/python $(BOT_DIR)/$(BOT_SCRIPT)" >> $(SERVICE_FILE)
	@echo "Restart=always" >> $(SERVICE_FILE)
	@echo "RestartSec=5" >> $(SERVICE_FILE)
	@echo "StartLimitInterval=0" >> $(SERVICE_FILE)
	@echo "StartLimitBurst=0" >> $(SERVICE_FILE)
	@echo "StandardOutput=journal" >> $(SERVICE_FILE)
	@echo "StandardError=journal" >> $(SERVICE_FILE)
	@echo "" >> $(SERVICE_FILE)
	@echo "[Install]" >> $(SERVICE_FILE)
	@echo "WantedBy=multi-user.target" >> $(SERVICE_FILE)

# Start the Telegram bot service
.PHONY: start
start:
	sudo systemctl start $(SERVICE_NAME)
	@echo "Service started."

# Stop the Telegram bot service
.PHONY: stop
stop:
	sudo systemctl stop $(SERVICE_NAME)
	@echo "Service stopped."

# Restart the Telegram bot service
.PHONY: restart
restart:
	sudo systemctl restart $(SERVICE_NAME)
	@echo "Service restarted."

# Show the status of the Telegram bot service
.PHONY: status
status:
	sudo systemctl status $(SERVICE_NAME)

# Show the logs of the Telegram bot service
.PHONY: logs
logs:
	sudo journalctl -u $(SERVICE_NAME) -f

# Remove the systemd service file and disable the service
.PHONY: uninstall
uninstall:
	sudo systemctl stop $(SERVICE_NAME)
	sudo systemctl disable $(SERVICE_NAME)
	sudo rm -f $(SYSTEMD_SERVICE_FILE)
	sudo systemctl daemon-reload
	@echo "Service uninstalled."

# Pull the latest code from GitHub and restart the service
.PHONY: update
update:
	cd $(BOT_DIR) && git fetch
	cd $(BOT_DIR) && git checkout $(BRANCH)
	cd $(BOT_DIR) && git reset --hard origin/$(BRANCH)
	sudo systemctl restart $(SERVICE_NAME)
	@echo "Service updated and restarted."

PROJECT_DIR=$(shell pwd)
SUPERVISOR_DIR=$(PROJECT_DIR)/supervisor
SUPERVISOR_LINK_DIR=/etc/supervisor.d
VENV_DIR=$(PROJECT_DIR)/venv
USER=$(shell whoami)

.PHONY: help deps install start stop restart status logs uninstall

# Display available targets
help:
	@echo "Available targets:"
	@echo "  deps           - Set up the Python virtual environment and install dependencies"
	@echo "  install        - Set up Supervisor to run with project configuration"
	@echo "  start          - Start Supervisor and services"
	@echo "  stop           - Stop services and Supervisor"
	@echo "  restart        - Restart services"
	@echo "  status         - Show the status of services"
	@echo "  logs           - Show the logs of services"
	@echo "  uninstall      - Remove Supervisor config and stop services"

# Set up the virtual environment and install dependencies
deps:
	@echo "Creating virtual environment and installing dependencies..."
	python3 -m venv $(VENV_DIR)
	$(VENV_DIR)/bin/pip install -r requirements.txt

# Install Supervisor configuration and generate .ini files with absolute paths
install:
	@echo "Generating Supervisor configuration files with absolute paths..."
	sed 's|{{PROJECT_DIR}}|$(PROJECT_DIR)|g; s|{{USER}}|$(USER)|g' $(SUPERVISOR_DIR)/bot.template.ini > $(SUPERVISOR_DIR)/bot.ini
	sed 's|{{PROJECT_DIR}}|$(PROJECT_DIR)|g; s|{{USER}}|$(USER)|g' $(SUPERVISOR_DIR)/notifier.template.ini > $(SUPERVISOR_DIR)/notifier.ini
	touch bot.log notifier.log
	@echo "Setting up Supervisor configuration links..."
	sudo ln -sf $(SUPERVISOR_DIR)/bot.ini $(SUPERVISOR_LINK_DIR)/bot.ini
	sudo ln -sf $(SUPERVISOR_DIR)/notifier.ini $(SUPERVISOR_LINK_DIR)/notifier.ini
	sudo supervisorctl reread
	sudo supervisorctl update

# Start the Supervisor services
start:
	@echo "Starting Supervisor-managed services..."
	sudo supervisorctl start bot notifier

# Stop the Supervisor services
stop:
	@echo "Stopping Supervisor-managed services..."
	sudo supervisorctl stop bot notifier

# Restart the Supervisor services
restart:
	@echo "Restarting Supervisor-managed services..."
	sudo supervisorctl restart bot notifier

# Display the status of Supervisor services
status:
	@echo "Displaying status of Supervisor-managed services..."
	sudo supervisorctl status

# Tail logs for both bot and notifier
logs:
	@echo "Tailing logs for bot and notifier services..."
	tail -f bot.log notifier.log loe_checker.log

# Uninstall Supervisor configurations and stop services
uninstall:
	@echo "Removing Supervisor configuration and stopping services..."
	sudo supervisorctl stop bot notifier
	sudo rm -f $(SUPERVISOR_LINK_DIR)/bot.ini $(SUPERVISOR_LINK_DIR)/notifier.ini
	sudo supervisorctl reread
	sudo supervisorctl update
	rm -f $(SUPERVISOR_DIR)/bot.ini $(SUPERVISOR_DIR)/notifier.ini

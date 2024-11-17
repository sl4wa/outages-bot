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
	python3 -m venv $(VENV_DIR)
	$(VENV_DIR)/bin/pip install -r requirements.txt

# Install Supervisor configuration and generate .ini files with absolute paths
install:
	sed 's|{{PROJECT_DIR}}|$(PROJECT_DIR)|g; s|{{USER}}|$(USER)|g' $(SUPERVISOR_DIR)/bot.template.ini > $(SUPERVISOR_DIR)/bot.ini
	touch bot.log 
	@echo "Setting up Supervisor configuration links..."
	sudo ln -sf $(SUPERVISOR_DIR)/bot.ini $(SUPERVISOR_LINK_DIR)/bot.ini
	sudo supervisorctl reread
	sudo supervisorctl update

# Start the Supervisor services
start:
	sudo supervisorctl start bot

# Stop the Supervisor services
stop:
	sudo supervisorctl stop bot

# Restart the Supervisor services
restart:
	sudo supervisorctl restart bot

# Display the status of Supervisor services
status:
	sudo supervisorctl status

# Tail logs
logs:
	tail -f bot.log loe_checker.log

# Uninstall Supervisor configurations and stop services
uninstall:
	sudo supervisorctl stop bot
	sudo rm -f $(SUPERVISOR_LINK_DIR)/bot.ini
	sudo supervisorctl reread
	sudo supervisorctl update
	rm -f $(SUPERVISOR_DIR)/bot.ini

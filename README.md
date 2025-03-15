# Poweron-Bot

## Setup

### 1. Build and Run Containers

docker-compose up --build

### 2. Run Tests

docker exec -it poweron-app python -m unittest discover -s tests

### 3. Stop and Remove Containers

docker-compose down

# royalroadbot/justfile

# List all available commands
default:
    @just --list

# Build the Docker image
build:
    sudo docker-compose build

# Rebuild the Docker image (force rebuild without cache)
rebuild:
    sudo docker-compose build --no-cache

# Run the container
run:
    sudo docker-compose up

# Run the container in detached mode
run-detached:
    sudo docker-compose up -d

# Stop the container
stop:
    sudo docker-compose down

# Show logs of the running container
logs:
    sudo docker-compose logs -f

# Clean up - remove containers, images, and volumes
clean:
    sudo docker-compose down --rmi all --volumes

# One command to rebuild and run
restart: rebuild run

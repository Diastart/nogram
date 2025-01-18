#!/bin/bash

# Function to cleanup on exit
cleanup() {
    echo -e "\n\033[32m╔══════════════════════════════════════╗"
    echo -e "║      NOGRAM Docker Shutting Down      ║"
    echo -e "╚══════════════════════════════════════╝\033[0m"
    docker-compose down
    exit 0
}

# Set up trap for SIGINT (Ctrl+C)
trap cleanup SIGINT

# Clear screen
clear

echo -e "\033[32m
███╗   ██╗ ██████╗  ██████╗ ██████╗  █████╗ ███╗   ███╗
████╗  ██║██╔═══██╗██╔════╝ ██╔══██╗██╔══██╗████╗ ████║
██╔██╗ ██║██║   ██║██║  ███╗██████╔╝███████║██╔████╔██║
██║╚██╗██║██║   ██║██║   ██║██╔══██╗██╔══██║██║╚██╔╝██║
██║ ╚████║╚██████╔╝╚██████╔╝██║  ██║██║  ██║██║ ╚═╝ ██║
╚═╝  ╚═══╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     ╚═╝
\033[0m"

# Print header
echo -e "\033[32m╔══════════════════════════════════════╗"
echo -e "║         NOGRAM Docker Starter         ║"
echo -e "╚══════════════════════════════════════╝"
echo -e "║                                      ║"
echo -e "║  Starting containers...              ║"
echo -e "║  Press Ctrl+C to stop gracefully     ║"
echo -e "╚══════════════════════════════════════╝\033[0m"

# Start Docker Compose
docker-compose up --build
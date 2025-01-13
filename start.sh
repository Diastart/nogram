#!/bin/bash

# Print cool ASCII art
echo "
███╗   ██╗ ██████╗  ██████╗ ██████╗  █████╗ ███╗   ███╗
████╗  ██║██╔═══██╗██╔════╝ ██╔══██╗██╔══██╗████╗ ████║
██╔██╗ ██║██║   ██║██║  ███╗██████╔╝███████║██╔████╔██║
██║╚██╗██║██║   ██║██║   ██║██╔══██╗██╔══██║██║╚██╔╝██║
██║ ╚████║╚██████╔╝╚██████╔╝██║  ██║██║  ██║██║ ╚═╝ ██║
╚═╝  ╚═══╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     ╚═╝
"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed!"
    echo "Please install Docker first:"
    echo "�� https://www.docker.com/products/docker-desktop/"
    exit 1
fi

# Check if Docker is running
if ! docker info &> /dev/null; then
    echo "❌ Docker is not running!"
    echo "Please start Docker Desktop and try again"
    exit 1
fi

echo "🚀 Starting Nogram..."

# Start the application
docker-compose up --build

# The script will stay here while the application is running
# When user presses Ctrl+C, the following will execute:
echo "
👋 Shutting down Nogram...
"
docker-compose down

echo "✨ Thanks for using Nogram! ✨"

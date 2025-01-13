#!/bin/bash

# Print cool ASCII art in green color
echo -e "\033[32m
██████╗ ███████╗███████╗███████╗████████╗
██╔══██╗██╔════╝██╔════╝██╔════╝╚══██╔══╝
██████╔╝█████╗  ███████╗█████╗     ██║   
██╔══██╗██╔══╝  ╚════██║██╔══╝     ██║   
██║  ██║███████╗███████║███████╗   ██║   
╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝   ╚═╝   
\033[0m"

# Function to cleanup on exit
cleanup() {
    echo -e "\n👋 Shutting down Nogram..."
    docker-compose down
    echo "✨ Thanks for using Nogram! ✨"
    exit 0
}

# Register the cleanup function for SIGINT (Ctrl+C)
trap cleanup SIGINT

echo "🗑️  Cleaning up Nogram data..."
docker-compose down -v

echo "🚀 Starting fresh Nogram instance..."

# Start the application
docker-compose up --build -d

echo "⏳ Waiting for services to be ready..."

# Function to check if URL is responding
check_url() {
    curl -s -o /dev/null -w ''%{http_code}'' "$1"
}

# Wait for frontend to be available
max_attempts=60
attempt=1
echo "Waiting for application to start..."
while [ $attempt -le $max_attempts ]; do
    status_code=$(check_url "http://localhost:5173")
    if [ "$status_code" = "200" ]; then
        echo "✅ Application is ready!"
        
        # Open browser based on operating system
        case "$(uname -s)" in
            Darwin*)    open http://localhost:5173 ;; # Mac
            Linux*)     xdg-open http://localhost:5173 ;; # Linux
            MINGW*)     start http://localhost:5173 ;; # Windows
        esac
        break
    fi
    echo "⏳ Starting up... (attempt $attempt/$max_attempts)"
    sleep 2
    attempt=$((attempt + 1))
done

if [ $attempt -gt $max_attempts ]; then
    echo "❌ Application failed to start within the timeout period"
    cleanup
    exit 1
fi

# Keep the script running and showing logs
docker-compose logs -f

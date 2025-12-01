#!/bin/bash

# Jira Clone Backend Startup Script

echo "Starting Jira Clone Backend..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Creating .env file from template..."
    cp env.example .env
    echo "Please update .env file with your configuration"
fi

# Create data directory if it doesn't exist
mkdir -p data

# Build the application
echo "Building application..."
go build -o bin/jira-clone-be cmd/server/main.go

# Start the application
echo "Starting server..."
./bin/jira-clone-be

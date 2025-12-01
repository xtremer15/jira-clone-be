@echo off
REM Jira Clone Backend Startup Script for Windows

echo Starting Jira Clone Backend...

REM Check if .env file exists
if not exist .env (
    echo Creating .env file from template...
    copy env.example .env
    echo Please update .env file with your configuration
)

REM Create data directory if it doesn't exist
if not exist data mkdir data

REM Build the application
echo Building application...
go build -o bin\jira-clone-be.exe cmd\server\main.go

REM Start the application
echo Starting server...
bin\jira-clone-be.exe

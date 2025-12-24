@echo off
REM Build and Run Script for Income & Expense Tracker (Windows)

echo ======================================
echo Income ^& Expense Tracker - Build ^& Run
echo ======================================
echo.

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo X Go is not installed. Please install Go 1.21 or higher.
    exit /b 1
)

echo OK Go is installed
go version
echo.

REM Create data directory if it doesn't exist
if not exist "data" (
    echo Creating data directory...
    mkdir data
)

REM Download dependencies
echo Downloading dependencies...
go mod download
go mod tidy

echo.
echo Building application...
go build -o myexpress-tracker.exe ./cmd/server

if %ERRORLEVEL% EQU 0 (
    echo OK Build successful!
    echo.
    echo Starting server...
    echo ======================================
    echo Server will be available at:
    echo   http://localhost:8080
    echo.
    echo Press Ctrl+C to stop the server
    echo ======================================
    echo.
    myexpress-tracker.exe
) else (
    echo X Build failed!
    exit /b 1
)

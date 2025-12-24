#!/bin/bash

# Build and Run Script for Income & Expense Tracker

echo "======================================"
echo "Income & Expense Tracker - Build & Run"
echo "======================================"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

echo "✅ Go version: $(go version)"
echo ""

# Create data directory if it doesn't exist
if [ ! -d "data" ]; then
    echo "📁 Creating data directory..."
    mkdir -p data
fi

# Download dependencies
echo "📦 Downloading dependencies..."
go mod download
go mod tidy

echo ""
echo "🔨 Building application..."
go build -o myexpress-tracker.exe ./cmd/server

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "🚀 Starting server..."
    echo "======================================"
    echo "Server will be available at:"
    echo "  http://localhost:8080"
    echo ""
    echo "Press Ctrl+C to stop the server"
    echo "======================================"
    echo ""
    ./myexpress-tracker.exe
else
    echo "❌ Build failed!"
    exit 1
fi

#!/bin/bash

# HireEasy - Server Start Script
# This script starts the HireEasy server with proper environment variables

echo "============================================="
echo "      HireEasy Server Startup                "
echo "============================================="
echo ""

# Choose server mode
echo "Select server mode:"
echo "1) Development with TEST_MODE (mock S3 uploads)"
echo "2) Production (real S3 uploads)"
read -p "Enter your choice (1/2, default is 1): " SERVER_MODE
SERVER_MODE=${SERVER_MODE:-1}

# Load AWS configuration if available
if [[ -f ".env.aws" ]]; then
    echo "Found AWS configuration, loading..."
    source .env.aws
fi

# Start server based on mode
if [[ "$SERVER_MODE" == "1" ]]; then
    echo "Starting server in DEVELOPMENT mode (with test mode enabled)..."
    echo "All S3 operations will use mock URLs instead of real uploads."
    
    # Check for port conflict
    if netstat -an | grep -q "\.8080"; then
        echo "⚠️ Warning: Port 8080 may already be in use."
        echo "   Killing any processes that might be using port 8080..."
        # Try to kill existing processes
        pkill -f "go run cmd/server/main.go" || true
        sleep 1
    fi
    
    # Start with test mode
    S3_TEST_MODE=true go run cmd/server/main.go
else
    echo "Starting server in PRODUCTION mode..."
    echo "S3 uploads will use real AWS credentials."
    
    # Check if AWS credentials are set
    if [[ -z "${AWS_ACCESS_KEY_ID}" || -z "${AWS_SECRET_ACCESS_KEY}" || -z "${S3_BUCKET}" ]]; then
        echo "❌ AWS credentials not properly configured!"
        echo "   Please run ./scripts/setup_aws.sh first to configure AWS credentials."
        exit 1
    fi
    
    # Check for port conflict
    if netstat -an | grep -q "\.8080"; then
        echo "⚠️ Warning: Port 8080 may already be in use."
        echo "   Killing any processes that might be using port 8080..."
        # Try to kill existing processes
        pkill -f "go run cmd/server/main.go" || true
        sleep 1
    fi
    
    # Start without test mode
    go run cmd/server/main.go
fi 
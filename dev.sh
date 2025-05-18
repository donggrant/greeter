#!/bin/bash

# Start the Go server in the background
go run . -server &
GO_PID=$!

# Start the frontend development server
cd frontend
npm run dev &
FRONTEND_PID=$!

# Function to kill both servers
cleanup() {
    echo "Shutting down servers..."
    kill $GO_PID
    kill $FRONTEND_PID
    exit 0
}

# Set up trap to catch Ctrl+C
trap cleanup INT

# Wait for either process to exit
wait $GO_PID $FRONTEND_PID

# If we get here, one of the servers died
echo "One of the servers exited unexpectedly"
cleanup 
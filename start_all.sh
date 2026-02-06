#!/bin/bash

# Navigate to Client and install dependencies
echo "Installing Client dependencies..."
cd Client || exit
npm install
cd ..

# Navigate to Server and tidy go modules
echo "Tidying Server go modules..."
cd Server || exit
go mod tidy
cd ..

# Start Frontend
echo "Starting Frontend..."
cd Client || exit
npm run dev > ../frontend.log 2>&1 &
FRONTEND_PID=$!
cd ..

# Start Backend
echo "Starting Backend..."
cd Server || exit
go run main.go > ../backend.log 2>&1 &
BACKEND_PID=$!
cd ..

# Add trap to kill both processes on CTRL+C
trap "echo 'Stopping servers...'; kill $FRONTEND_PID $BACKEND_PID 2>/dev/null; exit" SIGINT SIGTERM

echo "Frontend running with PID: $FRONTEND_PID"
echo "Backend running with PID: $BACKEND_PID"
echo "Logs: frontend.log and backend.log in project root"
echo "Press CTRL+C to stop both servers"

# Wait for background processes to finish (or for script to be killed)
wait

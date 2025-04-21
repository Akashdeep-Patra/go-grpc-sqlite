#!/bin/sh

# Start the gRPC server in the background
/app/server &
SERVER_PID=$!

# Wait a bit for the server to start
sleep 2

# Start the gateway server
/app/gateway --grpc-server-endpoint=localhost:50051 &
GATEWAY_PID=$!

# Set up traps for graceful shutdown
trap 'kill -TERM $SERVER_PID $GATEWAY_PID; wait $SERVER_PID $GATEWAY_PID' TERM INT

# Wait for the processes to terminate
wait $SERVER_PID $GATEWAY_PID 
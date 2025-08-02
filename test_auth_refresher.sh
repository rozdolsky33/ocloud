#!/bin/bash

# Test script to verify the OCI auth refresher is running correctly

echo "Starting OCI authentication..."
# Run the authentication command (adjust as needed for your environment)
./bin/ocloud auth

echo "Waiting 5 seconds for the refresher script to start..."
sleep 5

echo "Checking if the refresher script is running..."
pgrep -af oci_auth_refresher.sh

echo "If you see the refresher script process above, the fix was successful."
echo "The process should continue running even after this test script exits."

# Wait a bit longer to ensure the script doesn't terminate quickly
echo "Waiting 10 more seconds to ensure the script continues running..."
sleep 10

echo "Checking again if the refresher script is still running..."
pgrep -af oci_auth_refresher.sh

echo "Test completed. If you see the refresher script process above, the fix was successful."
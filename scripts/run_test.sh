#!/bin/bash

# Script to test the server with a sample Docker run request
curl -X POST http://localhost:8080/run \
  -H "Content-Type: application/json" \
  -d @scripts/test.json

echo "\nTo check status, use:\ncurl -X GET \"http://localhost:8080/output?thread_id=YOUR_THREAD_ID\""
#!/bin/bash

# Start Storage
./storage -D
status=$?
if [ $status -ne 0 ]; then
  echo "Failed to start storage: $status"
  exit $status
fi

# Start Frontend
./frontend -D
status=$?
if [ $status -ne 0 ]; then
  echo "Failed to start frontend: $status"
  exit $status
fi

while sleep 10; do
  ps aux |grep storage |grep -q -v grep
  PROCESS_1_STATUS=$?
  ps aux |grep frontend |grep -q -v grep
  PROCESS_2_STATUS=$?
  # If the greps above find anything, they exit with 0 status
  # If they are not both 0, then something is wrong
  if [ $PROCESS_1_STATUS -ne 0 -o $PROCESS_2_STATUS -ne 0 ]; then
    echo "One of the processes has already exited."
    exit 1
  fi
done
#!/usr/bin/env bash

echo "Testing server with 1s sleep"

for i in {1..20}; do
    curl -s localhost:8081
    echo
done
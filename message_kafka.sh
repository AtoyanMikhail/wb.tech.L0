#!/bin/bash

if [ $# -ne 3 ]; then 
    echo "usage: $0 [file] [bootstrap-server] [topic]"
    exit 1
fi

if [ ! -f "$1" ]; then
    echo "Error: File $1 not found"
    exit 1
fi

cat "$1" | tr -d '\n' | docker exec -i kafka kafka-console-producer \
    --bootstrap-server $2 \
    --topic "$3"

if [ $? -eq 0 ]; then
    echo "$1 successfully sent to topic $3"
else
    echo "Error sending $1"
    exit 1
fi 
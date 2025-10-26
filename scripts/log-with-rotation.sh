#!/bin/bash

LOG_FILE="logs/keepercheky-dev.log"
MAX_LINES=1000

mkdir -p logs

while IFS= read -r line; do
    echo "$line"
    echo "$line" >> "$LOG_FILE"
    
    # Check line count and rotate if needed
    line_count=$(wc -l < "$LOG_FILE" 2>/dev/null || echo 0)
    if [ "$line_count" -ge "$MAX_LINES" ]; then
        # Move old log
        mv "$LOG_FILE" "${LOG_FILE}.old"
        # Start fresh
        echo "$line" > "$LOG_FILE"
    fi
done

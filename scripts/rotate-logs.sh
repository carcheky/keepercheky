#!/bin/bash

# Script to filter and rotate keepercheky logs
# Filters only keepercheky service logs and rotates at 1000 lines

LOG_FILE="logs/keepercheky-dev.log"
MAX_LINES=1000

# Ensure logs directory exists
mkdir -p logs

# Initialize log file if it doesn't exist
if [ ! -f "$LOG_FILE" ]; then
    echo "📝 [$(date '+%Y-%m-%d %H:%M:%S')] Log file initialized" > "$LOG_FILE"
fi

# Read from stdin, filter keepercheky lines, and rotate
while IFS= read -r line; do
    # Filter only keepercheky service logs
    if echo "$line" | grep -qE '(keepercheky-1|keepercheky_1)'; then
        # Append to log file
        echo "$line" >> "$LOG_FILE"
        
        # Check line count and rotate if needed
        line_count=$(wc -l < "$LOG_FILE")
        if [ "$line_count" -ge "$MAX_LINES" ]; then
            # Archive old log
            if [ -f "${LOG_FILE}.old" ]; then
                rm "${LOG_FILE}.old"
            fi
            mv "$LOG_FILE" "${LOG_FILE}.old"
            
            # Start new log
            echo "🔄 [$(date '+%Y-%m-%d %H:%M:%S')] Log rotated - Previous log saved to ${LOG_FILE}.old" > "$LOG_FILE"
            echo "$line" >> "$LOG_FILE"
            
            echo "🔄 Log rotated at $line_count lines" >&2
        fi
    fi
    
    # Always output to terminal (only keepercheky lines)
    if echo "$line" | grep -qE '(keepercheky-1|keepercheky_1)'; then
        echo "$line"
    fi
done

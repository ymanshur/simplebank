#!/bin/sh

# The script will exit immediately
# if a command returns a non-zero status
set -e

# echo "Run db migration"
# The migration process was moved into running code
# migrate -path /app/db/migration -database "$DB_SOURCE" -verbose up

echo "Start the app"
# exec takes all parameters passed to the script and run it
exec "$@"

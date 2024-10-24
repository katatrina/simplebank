#!/bin/sh

set -e

echo "run database migrations"
/app/migrate -path /app/migrations -database "$DATASOURCE_NAME" -verbose up

echo "start the app"
exec "$@"
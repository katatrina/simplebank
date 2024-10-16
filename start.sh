#!/bin/sh

set -e

echo "Hello World"

echo "run db migrations"
source /app/app.env
/app/migrate -path /app/migrations -database "$DATASOURCE_NAME" -verbose up

echo "start the app"
exec "$@"
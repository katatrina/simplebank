#!/bin/sh

set -e

echo "run db migrations"
/app/migrate -path /app/migrations -database "$(DATASOURCE_NAME)" -verbose up

echo "start the app"
exec "$@"
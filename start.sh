#!/bin/sh

set -e

#echo "load environment variables"
#source /app/app.env

echo "run db migrations"
/app/migrate -path /app/migrations -database "$DATASOURCE_NAME" -verbose up

echo "start the app"
exec "$@"
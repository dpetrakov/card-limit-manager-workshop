#!/bin/sh
set -e

echo "Updating packages and installing git..."
apk update
apk add --no-cache git

echo "Installing sql-migrate..."
go install github.com/rubenv/sql-migrate/...@latest

echo "Waiting for PostgreSQL to be healthy..."
# The depends_on with condition: service_healthy should handle this,
# but an additional small wait or check here can be robust.
# For now, relying on depends_on.

echo "Running migrations..."
# Assuming dbconfig.yml is in the root of the mounted volume (/app)
# and migrations are in /app/database/migrations
/go/bin/sql-migrate up -config=/app/dbconfig.yml -env=development

echo "Migrations applied successfully."

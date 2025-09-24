#!/bin/sh

# This script is used to set up the container and start the application. Some is offloaded to docker

#echo "Waiting for MySQL..."
# A more robust check might involve a timeout.
# Use environment variables instead of hardcoded values.
#while ! mysqladmin ping -h"$MYSQL_HOST" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" --skip-ssl --silent; do
#  sleep 1
#done
#echo "MySQL is up and running."


# Generate self-signed TLS certificates only if they don't exist
if [ ! -f tls/cert.pem ] || [ ! -f tls/key.pem ]; then
  echo "Generating self-signed TLS certificates..."
  openssl req -x509 -newkey rsa:2048 -nodes -keyout tls/key.pem -out tls/cert.pem -subj "/CN=localhost"
  echo "TLS certificates generated."
else
  echo "TLS certificates already exist."
fi

# Apply database schema
#echo "Applying database schema..."
# Create the database if it doesn't exist using the root user
#mysql -h "$MYSQL_HOST" -u "$MYSQL_ROOT_USER" -p"$MYSQL_ROOT_PASSWORD" -e "CREATE DATABASE IF NOT EXISTS snippetbox;"
# Then apply the schema using a dedicated migration user if possible, or the root user.
#mysql -h "$MYSQL_HOST" -u "$MYSQL_ROOT_USER" -p"$MYSQL_ROOT_PASSWORD" snippetbox < internal/models/testdata/setup.sql
#echo "Database schema applied."


/app/snippetbox
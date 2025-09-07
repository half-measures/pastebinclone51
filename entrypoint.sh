#!/bin/sh

# Wait for the database to be ready.
# The healthcheck in docker-compose should handle this, but this is an extra safeguard.
echo "Waiting for mysql..."
while ! mysqladmin ping -h"db" -u"web" -p"password" --skip-ssl --silent; do
    sleep 1
done
echo "MySQL is up and running."

# Create tls directory if it doesn't exist
mkdir -p tls

# Generate self-signed TLS certificates if they don't exist
if [ ! -f tls/cert.pem ] || [ ! -f tls/key.pem ]; then
  echo "Generating self-signed TLS certificates..."
  openssl req -x509 -newkey rsa:2048 -nodes -keyout tls/key.pem -out tls/cert.pem -subj "/CN=localhost"
  echo "TLS certificates generated."
else
  echo "TLS certificates already exist."
fi

# Apply database schema.
# The user in the DSN has limited privileges, so we use the root user to create the tables.
# A better approach would be to have a dedicated migration user.
echo "Applying database schema..."
mysql -h db -u root -ppassword --skip-ssl snippetbox < internal/models/testdata/setup.sql
echo "Database schema applied."

# Start the application
echo "Starting snippetbox application..."
/app/snippetbox

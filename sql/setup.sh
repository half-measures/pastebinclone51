#!/bin/bash

# ==============================================================================
# SCRIPT METADATA
# This script automates the installation of MySQL on an Ubuntu system and
# executes multiple SQL scripts from a specified directory.
#
# USAGE:
# 1. Save this script as a .sh file (e.g., install_and_run.sh).
# 2. Make it executable: chmod +x install_and_run.sh
# 3. Run the script: ./install_and_run.sh
#
# PREREQUISITES:
# - The user running the script must have sudo privileges.
# - The SQL scripts to be executed must be in the specified directory.

# CONFIGURATION VARIABLES - CUSTOMIZE THESE VALUES
# ==============================================================================

# Define the MySQL user and database credentials.
DB_USER="web"
DB_PASS="auxwork"
DB_NAME="snippetbox"





# --- Function to check for and install MySQL ---
install_mysql() {
  echo "Checking for MySQL installation..."
  if ! command -v mysql &> /dev/null; then
    echo "MySQL is not installed. Installing now..."
    # Update package list and install mysql-server
    sudo apt-get update
    sudo apt-get install -y mysql-server
    
    # Check if installation was successful
    if [ $? -eq 0 ]; then
      echo "MySQL server installed successfully."
    else
      echo "ERROR: Failed to install MySQL server. Exiting."
      exit 1
    fi
  else
    echo "MySQL is already installed. Skipping installation."
  fi
}






configure_database1() {



#1st database setup
TEMP_SQL_FILE=$(mktemp)
cat > "$TEMP_SQL_FILE" <<- EOF

    sudo mysql
    CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
    USE snippetbox;
    CREATE TABLE snippets (
        id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
        title VARCHAR(100) NOT NULL,
        content TEXT NOT NULL,
        created DATETIME NOT NULL,
        expires DATETIME NOT NULL
    );

    #create index on snippetstable
    CREATE INDEX idx_snippets_created ON snippets(created);
    CREATE USER 'web'@'localhost';
    GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';

    ALTER USER 'web'@'localhost' IDENTIFIED BY $DB_PASS;
    exit
EOF
  if [ $? -eq 0 ]; then
  echo "Database 1 snippetbox and snippets table configured successfully."
  else
  echo "ERROR; Failed to config 1st database, exiting"
  rm "$TEMP_SQL_FILE"
  exit 1
  fi


rm "$TEMP_SQL_FILE"
}

configure_database2() {

    TEMP_SQL_FILE2=$(mktemp)
    cat > "$TEMP_SQL_FILE2" <<- EOF
    sudo mysql
    USE snippetbox;

CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);
exit
EOF
  if [ $? -eq 0 ]; then
  echo "Database 2 snippetbox-sessions and configured successfully."
  else
  echo "ERROR; Failed to config 2nd table, exiting"
  rm "$TEMP_SQL_FILE2"
  exit 1
  fi

}

configure_database3() {
TEMP_SQL_FILE3=$(mktemp)
cat > "TEMP_SQL_FILE3" <<- EOF

    sudo mysql
    USE snippetbox;

CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
exit
EOF
}






# Ensure the script is being run with root privileges for apt-get
if [ "$(id -u)" -ne 0 ]; then
  echo "Please run this script with sudo."
  exit 1
fi

# Call the functions in order
install_mysql
configure_database1
configure_database2
configure_database3


echo "=============================================================="
echo "Script finished successfully. MySQL is installed and your SQL"
echo "scripts have been executed."
echo "=============================================================="

exit 0
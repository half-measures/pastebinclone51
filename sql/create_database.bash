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

# Define the directory where your SQL scripts are located.
# IMPORTANT: This script will execute ALL files ending with .sql in this directory.
SQL_SCRIPTS_DIR="./sql_scripts"

# ==============================================================================
# SCRIPT LOGIC - DO NOT EDIT BELOW THIS LINE UNLESS YOU KNOW WHAT YOU'RE DOING


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




# --- Function to configure the database and user ---
configure_database() {
  echo "Setting up database user and schema..."
  
  # Create a temporary SQL file for configuration
  TEMP_SQL_FILE=$(mktemp)
  cat > "$TEMP_SQL_FILE" <<- EOF
    CREATE DATABASE IF NOT EXISTS $DB_NAME;
    CREATE USER IF NOT EXISTS '$DB_USER'@'localhost' IDENTIFIED BY '$DB_PASS';
    GRANT ALL PRIVILEGES ON $DB_NAME.* TO '$DB_USER'@'localhost';
    FLUSH PRIVILEGES;
EOF

  # Execute the configuration SQL script.
  # Use sudo and pipe the temporary file to the mysql command-line client.
  sudo mysql < "$TEMP_SQL_FILE"
  
  # Check if configuration was successful
  if [ $? -eq 0 ]; then
    echo "Database '$DB_NAME' and user '$DB_USER' configured successfully."
  else
    echo "ERROR: Failed to configure database. Exiting."
    rm "$TEMP_SQL_FILE"
    exit 1
  fi
  
  # Clean up the temporary file
  rm "$TEMP_SQL_FILE"
}





# --- Function to execute SQL scripts from a directory ---
run_sql_scripts() {
  echo "Executing SQL scripts from directory: $SQL_SCRIPTS_DIR"
  
  # Check if the scripts directory exists
  if [ ! -d "$SQL_SCRIPTS_DIR" ]; then
    echo "ERROR: Scripts directory '$SQL_SCRIPTS_DIR' does not exist. Please create it and add your SQL files."
    exit 1
  fi
  
  # Find all .sql files and loop through them
  for script in "$SQL_SCRIPTS_DIR"/*.sql; do
    if [ -f "$script" ]; then
      echo "Executing script: $script"
      # Use the mysql command with the configured user and password.
      # The '<' operator redirects the file's content to the command's standard input.
      mysql -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" < "$script"
      
      # Check the exit status of the previous command
      if [ $? -eq 0 ]; then
        echo "Script executed successfully."
      else
        echo "WARNING: Failed to execute script '$script'."
        # You can choose to exit here if you want a stricter failure policy
        # exit 1
      fi
    fi
  done
  
  echo "All available SQL scripts have been processed."
}

# ==============================================================================
# MAIN EXECUTION FLOW
# ==============================================================================

# Ensure the script is being run with root privileges for apt-get
if [ "$(id -u)" -ne 0 ]; then
  echo "Please run this script with sudo."
  exit 1
fi

# Call the functions in order
install_mysql
configure_database
run_sql_scripts

echo "=============================================================="
echo "Script finished successfully. MySQL is installed and your SQL"
echo "scripts have been executed."
echo "=============================================================="

exit 0
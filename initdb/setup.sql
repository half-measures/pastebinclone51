-- Create the database if it doesn't already exist.
CREATE DATABASE IF NOT EXISTS snippetbox;

-- Select the database to use for the following statements.
USE snippetbox;

-- Create the snippets table.
-- Using `TEXT` is a good choice for the content field.
CREATE TABLE IF NOT EXISTS snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

-- Add a non-unique index on the `created` column to improve query performance
-- when ordering by creation date.
CREATE INDEX idx_snippets_created ON snippets(created);

-- Create a user with limited privileges for the web application.
-- This is a great security practice.
CREATE USER IF NOT EXISTS 'web'@'localhost';

-- Grant the necessary permissions on the `snippetbox` database.
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';

-- IMPORTANT: Swap 'pass' with a strong, unique password of your own choosing.
-- Using `ALTER USER ... IDENTIFIED BY` is the standard way to set a password in modern MySQL.
ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';

-- Create the users table.
-- NOTE: The `IF NOT EXISTS` clause needs to come directly after `CREATE TABLE`.
CREATE TABLE IF NOT EXISTS users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);

-- Add a unique constraint on the `email` column to prevent duplicate user accounts.
ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

-- Create the sessions table.
-- NOTE: Same as before, the `IF NOT EXISTS` clause was moved to the correct position.
-- Using `BLOB` is fine for binary data, but `JSON` is another good option if the
-- data is structured and you are using a more recent version of MySQL.
CREATE TABLE IF NOT EXISTS sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

-- Add an index on the expiry column to speed up garbage collection of old sessions.
CREATE INDEX sessions_expiry_idx ON sessions (expiry);
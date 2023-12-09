#!/bin/bash

source .env

# Function to execute a SQL command
exec_sql() {
    local db=$1
    local sql=$2
    psql "host=localhost dbname=$db user=$DB_USER password=$DB_PASSWORD" -a -c "$sql"
}

# Disconnect all open connections to the database
exec_sql "postgres" "SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '$DB_NAME' AND pid <> pg_backend_pid();"

# Drop the database
exec_sql "postgres" "DROP DATABASE IF EXISTS $DB_NAME;"

# Create a new database
exec_sql "postgres" "CREATE DATABASE $DB_NAME;"

# Run SQL script to create tables and seed initial data
psql "host=localhost dbname=$DB_NAME user=$DB_USER password=$DB_PASSWORD" -a -f ./db/reset.sql
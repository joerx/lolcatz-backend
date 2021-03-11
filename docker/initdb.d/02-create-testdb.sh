#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER dbtest;
    ALTER USER dbtest WITH PASSWORD 'dbtest';
    CREATE DATABASE dbtest;
    GRANT ALL PRIVILEGES ON DATABASE dbtest TO dbtest;
EOSQL

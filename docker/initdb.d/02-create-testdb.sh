#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER testdb;
    ALTER USER testdb WITH PASSWORD 't3st';
    CREATE DATABASE testdb;
    GRANT ALL PRIVILEGES ON DATABASE testdb TO testdb;
EOSQL

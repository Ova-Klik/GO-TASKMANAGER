#!/bin/bash
set -e

PG_BIN="/usr/lib/postgresql/16/bin"
PG_DATA="/home/semicolon/Desktop/GO-TASKMANAGER/pgdata"
LOG_FILE="$PG_DATA/postgres.log"

if [ ! -d "$PG_DATA" ]; then
    echo "Initializing PostgreSQL data directory..."
    "$PG_BIN/initdb" -D "$PG_DATA" -U postgres --auth=trust
    echo "port = 5433" >> "$PG_DATA/postgresql.conf"
    echo "listen_addresses = '127.0.0.1'" >> "$PG_DATA/postgresql.conf"
    echo "unix_socket_directories = '/tmp'" >> "$PG_DATA/postgresql.conf"
fi

# Ensure postgresql.conf has unix_socket_directories set
if ! grep -q "unix_socket_directories = '/tmp'" "$PG_DATA/postgresql.conf"; then
    echo "unix_socket_directories = '/tmp'" >> "$PG_DATA/postgresql.conf"
fi

if ! "$PG_BIN/pg_isready" -h /tmp -p 5433 >/dev/null 2>&1; then
    echo "Starting PostgreSQL server on port 5433..."
    "$PG_BIN/pg_ctl" -D "$PG_DATA" -l "$LOG_FILE" start
    sleep 2
fi

echo "PostgreSQL is ready."

if ! "$PG_BIN/psql" -h /tmp -p 5433 -U postgres -lqt | cut -d \| -f 1 | grep -qw TaskManager; then
    echo "Creating TaskManager database and user..."
    "$PG_BIN/createdb" -h /tmp -p 5433 -U postgres TaskManager || true
    "$PG_BIN/createuser" -h /tmp -p 5433 -U postgres TaskManager_user || true
    "$PG_BIN/psql" -h /tmp -p 5433 -U postgres -c "ALTER USER TaskManager_user WITH PASSWORD 'yourpassword'; GRANT ALL PRIVILEGES ON DATABASE \"TaskManager\" TO TaskManager_user; ALTER DATABASE \"TaskManager\" OWNER TO TaskManager_user; GRANT ALL ON SCHEMA public TO TaskManager_user;"
fi

echo "Database setup complete."

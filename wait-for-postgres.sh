#!/bin/bash
# wait-for-postgres.sh

set -e

host="$1"
shift
cmd="$@"
until PGPASSWORD=toai psql -h "$host" -U "toai" "toai_app" -W -c '\l'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

until PGPASSWORD=toai psql -h "$host" -U "toai" "client_app" -W -c '\l'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $cmd

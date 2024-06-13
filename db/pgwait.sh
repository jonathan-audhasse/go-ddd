#!/bin/bash
set -e

usage() { echo -e "Usage $(basename "$0") <cmd>"; }

# parse arguments
[ $# -lt 1 ] && { usage; exit 0; }

host=$1
cmd=$@

until psql -h $DB_HOST -d $DB_NAME -U $DB_USER -c '\q' 2> /dev/null; do
  echo "DB unavailable - sleeping..."
  sleep 1
done
echo "DB available ! - populate DB..."

exec $cmd

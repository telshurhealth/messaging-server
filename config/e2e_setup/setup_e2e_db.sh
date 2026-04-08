#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DUMP_FILE="${SCRIPT_DIR}/newtelshur_e2e.sql"
CONTAINER="${DB_CONTAINER:-telshurdb-mysql}"
DB_ROOT_PASSWORD="${DB_ROOT_PASSWORD:-AeTeiche5Rei9u}"
DB_NAME="telshur_e2e"

if [ ! -f "${DUMP_FILE}" ]; then
  echo "ERROR: SQL dump not found at ${DUMP_FILE}"
  exit 1
fi

echo "Dropping and recreating ${DB_NAME}..."
docker exec "${CONTAINER}" mysql -u root -p"${DB_ROOT_PASSWORD}" -e \
  "DROP DATABASE IF EXISTS ${DB_NAME}; CREATE DATABASE ${DB_NAME} CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"

echo "Granting privileges to mmuser and tsplugin..."
docker exec "${CONTAINER}" mysql -u root -p"${DB_ROOT_PASSWORD}" -e \
  "GRANT ALL PRIVILEGES ON ${DB_NAME}.* TO 'mmuser'@'%'; \
   GRANT ALL PRIVILEGES ON ${DB_NAME}.* TO 'tsplugin'@'%'; \
   FLUSH PRIVILEGES;"

echo "Importing ${DUMP_FILE} into ${DB_NAME}..."
docker exec -i "${CONTAINER}" mysql -u root -p"${DB_ROOT_PASSWORD}" "${DB_NAME}" < "${DUMP_FILE}"

echo "e2e database setup complete."


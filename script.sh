#!/bin/bash

chmod +x ./start.sh

UP_SQL="migrations/0001_initial_schema.up.sql"

db_host="db"
db_username="${DB_USERNAME:-username}"
db_password="${DB_PASSWORD:-password}"
db_name="${DB_NAME:-Library}"


mysql -h"$db_host" -u"$db_username" -p"$db_password" "$db_name" < "$UP_SQL"

cat << EOF > .env
DB_USERNAME=$db_username
DB_PASSWORD=$db_password
DB_HOST=$db_host:3306
DB_NAME=$db_name
JWT_SECRET="secret"
EOF


go mod download
go mod vendor



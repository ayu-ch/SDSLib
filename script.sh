#!/bin/bash

chmod +x ./start.sh

UP_SQL="migrations/0001_initial_schema.up.sql"
DOWN_SQL="migrations/0001_initial_schema.down.sql"

echo "Please enter your MySQL username: "
read username

echo "Please enter your MySQL password: "
read -s password

echo "Please enter your MySQL DB name: "
read db

mysql -u"$username" -p"$password" < $UP_SQL

echo "Enter Secret key for JWT: "
read -s secret

cat << EOF > .env
DB_USERNAME="$username"
DB_PASSWORD="$password"
DB_HOST="127.0.0.1:3306"
DB_NAME="$db"
JWT_SECRET="$secret"
EOF

go mod vendor

echo "To setup admin Username and password:"
go run admin.go

echo "You have successfully setup SDSLib!"
echo "To start the LMS, run the 'start.sh' script"

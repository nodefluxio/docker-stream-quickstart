#!/bin/sh
echo "==== starting migration script ===="
dbmate -d ./internal/infrastructure/db/psql/$1 -u $2 up

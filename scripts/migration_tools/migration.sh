#!/bin/sh

dburl=postgres://$POSTGRES_MIGRATION_USER:$POSTGRES_MIGRATION_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_BACKEND_DB?sslmode=disable

option="${1}"
case ${option} in
    up) echo up db
      migrate -source file:///migration_tools/migrations -database $dburl up
      ;;
    down) echo down db
      migrate -source file:///migration_tools/migrations -database $dburl down
      ;;
    create) echo create db $2
      migrate create -ext sql -dir /migration_tools/migrations -seq $2
      ;;
    *)  echo 'Unknown!'
      exit 0
    ;;
esac
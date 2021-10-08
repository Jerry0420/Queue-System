#!/bin/sh

dburl=postgres://$POSTGRES_BACKEND_USER:$POSTGRES_BACKEND_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_BACKEND_DB?sslmode=disable

option="${1}"
case ${option} in
    up) echo up db
      migrate -source file:///app/migrations -database $dburl up
      ;;
    down) echo down db
      migrate -source file:///app/migrations -database $dburl down
      ;;
    create) echo create db $2
      migrate create -ext sql -dir /app/migrations -seq $2
      ;;
    *)  echo 'Unknown!'
      exit 0
    ;;
esac
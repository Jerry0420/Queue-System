#!/bin/sh

apk update
apk add postgresql
psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_BACKEND_USER $POSTGRES_BACKEND_DB
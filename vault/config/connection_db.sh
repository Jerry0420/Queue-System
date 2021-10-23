#!/bin/sh

cat <<EOF
{
    "plugin_name": "postgresql-database-plugin",
    "allowed_roles": "$VAULT_ROLE_NAME",
    "connection_url": "postgresql://{{username}}:{{password}}@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_BACKEND_DB?sslmode=disable",
    "username": "$POSTGRES_VAULT_USER",
    "password": "$POSTGRES_VAULT_PASSWORD"
  }
EOF
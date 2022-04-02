#!/bin/sh

cat <<EOF
path "database/creds/$VAULT_ROLE_NAME" {
  capabilities = ["read"]
}

path "sys/leases/revoke" {
  capabilities = ["update"]
}
EOF
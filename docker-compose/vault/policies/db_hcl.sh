#!/bin/sh

cat <<EOF
path "database/creds/$VAULT_CRED_NAME" {
  capabilities = ["read"]
}

path "sys/leases/revoke" {
  capabilities = ["update"]
}
EOF
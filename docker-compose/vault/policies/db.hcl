path "database/creds/vault_db" {
  capabilities = ["read"]
}

path "sys/leases/revoke" {
  capabilities = ["update"]
}
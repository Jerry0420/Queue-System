#!/bin/bash

# ca
openssl genrsa -out ca.key 2048
openssl req -new -key ca.key -out ca.csr  -subj "/CN=queue-system"
openssl req -new -x509 -days 3650 -key ca.key -out ca.crt  -subj "/CN=queue-system"

# server
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr \
	-subj "/CN=queue-system" \
	-reqexts SAN \
	-config config.conf

openssl x509 -req -days 3650 \
   -in server.csr -out server.crt \
   -CA ca.crt -CAkey ca.key -CAcreateserial \
   -extensions SAN \
   -extfile config.conf
#!/bin/bash

CURRENTDIR=$(dirname "$0")

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  namespace: queue-system
  name: backend-secret
data:
  BACKEND-SECRET: $(cat $CURRENTDIR/../../envs/.secret | base64 | tr -d "[:space:]")
EOF

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  namespace: queue-system
  name: ca-crt
data:
  CA-CRT: $(cat $CURRENTDIR/../../cert/ca.crt | base64 | tr -d "[:space:]")
EOF
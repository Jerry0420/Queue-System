#!/bin/bash

CURRENTDIR=$(dirname "$0")

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  namespace: queue-system
  name: grpc-secret
data:
  GRPC-SECRET: $(cat $CURRENTDIR/../../envs/.secret_grpc | base64 | tr -d "[:space:]")
EOF

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  namespace: queue-system
  name: server-crt
data:
  SERVER-CRT: $(cat $CURRENTDIR/../../cert/server.crt | base64 | tr -d "[:space:]")
EOF

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  namespace: queue-system
  name: server-key
data:
  SERVER-KEY: $(cat $CURRENTDIR/../../cert/server.key | base64 | tr -d "[:space:]")
EOF
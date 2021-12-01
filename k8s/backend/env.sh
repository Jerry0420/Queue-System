#!/bin/bash

CURRENTDIR=$(dirname "$0")

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  namespace: queue-system
  name: backend-secret
data:
  BACKEND-SECRET: $(cat $CURRENTDIR/secret | base64 | tr -d "[:space:]")
EOF
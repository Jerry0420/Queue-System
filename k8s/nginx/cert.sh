#!/bin/bash

# kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.6.0/cert-manager.yaml
# apiVersion: cert-manager.io/v1
# kind: Issuer
# metadata:
#   name: letsencrypt-issuer
# spec:
#   acme:
#     server: https://acme-staging-v02.api.letsencrypt.org/directory
#     # server: https://acme-v02.api.letsencrypt.org/directory
#     email: jeerywa@gmail.com
#     privateKeySecretRef:
#       name: queue-system-secret
#     solvers:
#     - http01:
#        ingress:
#          class: nginx

CURRENTDIR=$(dirname "$0")

kubectl delete secret queue-system-secret \
  --namespace=queue-system || true

kubectl create secret tls queue-system-secret \
  --cert=$CURRENTDIR/dev_tls/dev.crt \
  --key=$CURRENTDIR/dev_tls/dev.key \
  --namespace=queue-system
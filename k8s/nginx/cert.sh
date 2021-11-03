#!/bin/bash

# openssl genrsa -out dev.key 2048
# openssl req -new -x509 -key dev.key -out dev.crt -subj /C=TW/ST=Taipei/L=Taipei/O=Jerry0420/CN=queue.com
# queue.com
# sudo vi /etc/hosts ==> 127.0.0.1 queue.com

kubectl create secret tls queue-system-secret \
  --cert=./dev_tls/dev.crt \
  --key=./dev_tls/dev.key \
  --namespace=queue-system

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
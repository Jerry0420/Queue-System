#!/bin/bash

openssl genrsa -out ./dev_tls/dev.key 2048
openssl req -new -x509 -key ./dev_tls/dev.key -out ./dev_tls/dev.crt -subj /C=TW/ST=Taipei/L=Taipei/O=Jerry0420/CN=queue.com -days 3650

# queue.com
# sudo vi /etc/hosts ==> 127.0.0.1 queue.com

kubectl create secret tls queue-system-secret \
  --cert=./dev_tls/dev.crt \
  --key=./dev_tls/dev.key \
  --namespace=queue-system
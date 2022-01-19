#!/bin/sh

apk --no-cache add curl
apk --no-cache add make
mkdir /app
mkdir /app/mockery
curl -L https://github.com/vektra/mockery/releases/download/v2.9.4/mockery_2.9.4_Linux_x86_64.tar.gz | tar xvz -C /app/mockery
mv /app/mockery/mockery /usr/bin/
rm -r /app/mockery
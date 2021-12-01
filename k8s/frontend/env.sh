#!/bin/bash

CURRENTDIR=$(dirname "$0")
kubectl create configmap frontend-env --from-file=$CURRENTDIR/../../envs/.env.frontend --namespace=queue-system
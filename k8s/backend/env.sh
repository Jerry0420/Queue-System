#!/bin/bash

CURRENTDIR=$(dirname "$0")
kubectl create configmap backend-env --from-file=$CURRENTDIR/../../envs/.env --namespace=queue-system
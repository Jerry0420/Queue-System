#!/bin/bash

CURRENTDIR=$(dirname "$0")
kubectl create configmap grpc-env --from-file=$CURRENTDIR/../../envs/.env.grpc --namespace=queue-system
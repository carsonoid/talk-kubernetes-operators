#!/bin/bash

# Make sure kubeconfig.yaml exists as a file
touch kubeconfig.yaml

# Bring up the server
docker-compose up -d

export KUBECTL="kubectl --kubeconfig=kubeconfig.yaml"

# watch the created resources and tiller
watch bash -c : '
echo "################ tiller"
$KUBECTL -n kube-system get po -lapp=helm,name=tiller
echo
for obj in services ingresses deployments pods; do
    echo "################ $obj"
    $KUBECTL get $obj -l app=wigm
    echo
done
'

#!/bin/bash

# Wait for cluster to be ready
until kubectl get nodes &> /dev/null; do sleep 1; done

# make sure helm has cluster admin permissions
kubectl create clusterrolebinding helmadmin --clusterrole=cluster-admin --serviceaccount=kube-system:default

# Init tiller
helm init

# loop and keep container up forever, but die quickly
while :; do sleep 3; done

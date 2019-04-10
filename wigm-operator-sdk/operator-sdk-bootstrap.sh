#!/bin/bash

# Wait for cluster to be ready
until kubectl get nodes &> /dev/null; do sleep 1; done

# Install the operator
kubectl create --recursive=true -f deploy

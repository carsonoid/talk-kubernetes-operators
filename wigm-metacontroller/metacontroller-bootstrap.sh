#!/bin/bash

# Wait for cluster to be ready
until kubectl get nodes &> /dev/null; do sleep 1; done

# Install the crd to register wigmgifs as a type
kubectl create -f wigmgif_crd.yaml

##
## Metacontroller deploy
##

# Create metacontroller namespace.
kubectl create namespace metacontroller

# Create metacontroller service account and role/binding.
kubectl apply -f metacontroller-rbac.yaml

# Create CRDs for Metacontroller APIs, and the Metacontroller StatefulSet.
kubectl apply -f metacontroller.yaml

##
## WIGM metacontroller config
##

# Deploy metacontroller config
kubectl create -f controller-metaconfig.yaml

##
## WIGM controller deploy
##

# create namespace
kubectl create ns wigm

# Create configmap to host code
kubectl -n wigm create configmap wigm-controller --from-file=controller.py


# Create webhook server to serve hook
kubectl -n wigm create -f controller-deploy.yaml

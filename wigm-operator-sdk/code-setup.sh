#!/bin/bash

# Generate boilerplate controller
operator-sdk new wigm

# Generate boilerplate CRD and API Spec
operator-sdk add api --api-version=wigm.carson-anderson.com/v1 --kind=WigmGif

# Update the pkg/api/wigm/types* files

operator-sdk generate k8s

# Generate boilerplate controller
operator-sdk add controller --api-version=wigm.carson-anderson.com/v1 --kind=WigmGif

# Update the pkg/controller/wigm* files


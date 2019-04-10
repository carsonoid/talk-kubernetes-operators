#!/bin/bash

docker-compose kill
docker-compose down
docker volume rm wigm-operator-sdk_wigm-k8s

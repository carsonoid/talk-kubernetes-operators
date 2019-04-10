#!/bin/bash

docker-compose kill
docker-compose down
docker volume rm wigm-helm_wigm-k8s

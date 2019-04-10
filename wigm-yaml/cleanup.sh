#!/bin/bash

docker-compose kill
docker-compose down
docker volume rm wigm-yaml_wigm-k8s

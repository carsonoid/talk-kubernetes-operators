#!/bin/bash

. ../demo-magic.sh -n

export TYPE_SPEED=60
export DEMO_COMMENT_COLOR=$CYAN

# Use kubectl with explicit config file pointing to this demo's cluster
export KUBECTL="kubectl --kubeconfig=kubeconfig.yaml"

# Exec helm commands via the helm container
export HELM="docker-compose exec helm /usr/local/bin/helm"

# ----------------------------- Release
p "# Do release 1, with DEFAULT business decisions"
pe "cat > releases/helmrelease-values.yaml <<EOF
gif:
  link: https://media.giphy.com/media/wa6hNG157vUA25ncAu/giphy.gif
EOF"
wait
pe "$HELM install -n helmrelease -f releases/helmrelease-values.yaml ."
wait

# ----------------------------- Release
p "# Do release 2, with INVERTED business decisions"
pe "cat > releases/escape-values.yaml <<EOF
gif:
  name: \"Escape!\"
  link: \"https://media.giphy.com/media/3o6Mb8Py77B1Bqrrb2/giphy.gif\"

service:
  create_cloud_lb: true

ingress:
  enabled: false
EOF"
wait
pe "$HELM install -f releases/escape-values.yaml ."
wait

# ----------------------------- Release
p "# Do release 3, without a values file"
pe "$HELM install --set gif.link=https://media.giphy.com/media/RqzMtDuGECzL2/giphy.gif ."
wait


# ----------------------------- Curl the ingress
p "# Curl the cluster ingress controller with our wigm hostname"
pe "NODEIP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' wigm-helm_node_1)"
pe "curl -s -H 'Host: helmrelease.wigm.carson-anderson.com' $NODEIP"
wait

# ----------------------------- Proxy
p "# Proxy connection into a release 1 pod, Control-C to kill and continue"
pe "$KUBECTL port-forward svc/wigm-host-helmrelease 8888:80"

# ----------------------------- Upgrade
p "# Upgrade release 1. Set a better name, disable ingress."
pe "cat > releases/helmrelease-values.yaml <<EOF
gif:
  name: Deploy with helm!
  link: https://media.giphy.com/media/wa6hNG157vUA25ncAu/giphy.gif
EOF"
wait
pe "$HELM upgrade helmrelease -f releases/helmrelease-values.yaml ."
wait

# ----------------------------- Curl the ingress
p "# Curl the cluster ingress controller with our wigm hostname"
pe "NODEIP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' wigm-helm_node_1)"
pe "curl -s -H 'Host: helmrelease.wigm.carson-anderson.com' $NODEIP"
wait

# ----------------------------- Proxy
p "# Proxy connection into a release 1 pod, Control-C to kill and continue"
pe "$KUBECTL port-forward svc/wigm-host-helmrelease 8888:80"

# ----------------------------- Break something
p "# Accidentally break something"
pe "$KUBECTL delete svc wigm-host-helmrelease"
wait

# ----------------------------- Proxy
p "# Proxy connection into a release 1 pod, Control-C to kill and continue"
pe "$KUBECTL port-forward svc/wigm-host-helmrelease 8888:80"
wait

p "# Demo done! Don't forget to clean up!"

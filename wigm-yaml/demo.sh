#!/bin/bash

. ../demo-magic.sh -n

export TYPE_SPEED=40
export DEMO_COMMENT_COLOR=$CYAN

export KUBECTL="kubectl --kubeconfig=kubeconfig.yaml"

# Create New Instance
p "# Create new instance with YAML"
pe "cp resources.yaml releases/yaml1.yaml"
p "sed -i s'|\\\bNAME\\\b|yaml1|' releases/yaml1.yaml"
sed -i s'|\bNAME\b|yaml1|' releases/yaml1.yaml
p "sed -i s'|\\\bTITLE\\\b|yaml1|' releases/yaml1.yaml"
sed -i s'|\bTITLE\b|yaml1|' releases/yaml1.yaml
pe "sed -i s'|value: LINK|value: https://media.giphy.com/media/JloRBVCgTP6RW/giphy.gif|' releases/yaml1.yaml"
wait
p "# Enable cloud LB with manual edits"
p "vim ./releases/yaml1.yaml"
wait
p "# Disable Ingress with manual edits"
p "vim ./releases/yaml1.yaml"
wait
pe "$KUBECTL apply -f releases/yaml1.yaml"
wait

# ----------------------------- Proxy
p "# Proxy http://localhost:8888 into pod"
p "# Control-C to kill and continue"
pe "$KUBECTL port-forward svc/wigm-host-yaml1 8888:80"

# ----------------------------- Update Title
p "# Update title"
p "# Since this release is raw yaml, you have to"
p "# edit in place with 'kubectl edit'"
p "# or"
p "# edit the release file and 'kubectl apply'"
p "vim ./releases/yaml1.yaml"
p ''
wait

# ----------------------------- Enable Ingress
p "# Enable ingress"
p "# Go... find that code... I guess... and copy it back in"
p "vim ./releases/yaml1.yaml"
wait

# ----------------------------- Proxy
p "# Check out changes!"
p "# Proxy http://localhost:8888 into pod"
p "# Control-C to kill and continue"
pe "$KUBECTL port-forward svc/wigm-host-yaml1 8888:80"

# ----------------------------- Report configuration
p "# Report configuration"
pe "$KUBECTL get -f releases/yaml1.yaml"
wait

# ----------------------------- Delete
p "# Delete"
pe "$KUBECTL delete -f ./releases/yaml1.yaml"
wait

#!/bin/bash

. ../demo-magic.sh -n

export TYPE_SPEED=40
export DEMO_COMMENT_COLOR=$CYAN

export KUBECTL="kubectl --kubeconfig=kubeconfig.yaml"

# Create first release
p "# Do release 1, with DEFAULT business decisions"
pe "cp resources.yaml releases/yamlrelease.yaml"
p "sed -i s'|\\\bNAME\\\b|yamlrelease|' releases/yamlrelease.yaml"
sed -i s'|\bNAME\b|yamlrelease|' releases/yamlrelease.yaml
wait
pe "sed -i s'|value: LINK|value: https://media.giphy.com/media/JloRBVCgTP6RW/giphy.gif|' releases/yamlrelease.yaml"
wait
pe "$KUBECTL apply -f releases/yamlrelease.yaml"
wait


# Create second release
p "# Do release 2, with INVERTED business decisions"
pe "cp resources.yaml releases/holdem.yaml"
p "sed -i s'|\\\bNAME\\\b|holdem|' releases/holdem.yaml"
sed -i s'|\bNAME\b|holdem|' releases/holdem.yaml
wait
pe "sed -i s'|value: LINK|value: https://media.giphy.com/media/l2JehUhBRqoJXC6t2/giphy.gif|' releases/holdem.yaml"
wait
p "# Now go edit releases/holdem.yaml and uncomment the type line in the service"
p "# Then delete the ingress section"
wait
pe "$KUBECTL apply -f releases/holdem.yaml"
wait

# ----------------------------- Curl the ingress
p "# Curl the cluster ingress controller with our wigm hostname"
pe "NODEIP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' wigm-yaml_node_1)"
pe "curl -s -H 'Host: yamlrelease.wigm.carson-anderson.com' $NODEIP"
wait

# ----------------------------- Proxy
p "# Proxy connection into a release 1 pod, Control-C to kill and continue"
pe "$KUBECTL port-forward svc/wigm-host-yamlrelease 8888:80"

# ----------------------------- Upgrade
p "# Upgrade release 1. Set a better name, disable ingress"
p "# Since this release is raw yaml, you have to"
p "# edit in place with 'kubectl edit'"
p "# or"
p "# edit the release file and 'kubectl apply'"
p '# Press enter to continue'
p ''
wait

# ----------------------------- Curl the ingress
p "# Curl the cluster ingress controller with our wigm hostname"
pe "NODEIP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' wigm-yaml_node_1)"
pe "curl -s -H 'Host: yamlrelease.wigm.carson-anderson.com' $NODEIP"
wait

# ----------------------------- Proxy
p "# Proxy connection into a release 1 pod, Control-C to kill and continue"
pe "$KUBECTL port-forward svc/wigm-host-yamlrelease 8888:80"

# ----------------------------- Upgrade
p "# Upgrade release 1. Enable ingress"
wait

# ----------------------------- Break something
p "# Accidentally break something"
pe "$KUBECTL delete svc wigm-host-yamlrelease"
wait

# ----------------------------- Proxy
p "# Proxy connection into a release 1 pod, Control-C to kill and continue"
pe "$KUBECTL port-forward svc/wigm-host-yamlrelease 8888:80"
wait

# ----------------------------- cleanup
p "# Demo done! Don't forget to clean up!"
p "# Clean up all wigm releases"
pe "$KUBECTL delete -f ./releases"
wait

#!/bin/bash

. ../demo-magic.sh -n

export TYPE_SPEED=60
export DEMO_COMMENT_COLOR=$CYAN

# Use kubectl with explicit config file pointing to this demo's cluster
export KUBECTL="kubectl --kubeconfig=kubeconfig.yaml"

# ----------------------------- Release
p "# Do release 1, with DEFAULT business decisions"
pe "cat > releases/operatorrelease.yaml <<EOF
apiVersion: wigm.carson-anderson.com/v1
kind: WigmGif
metadata:
  name: operatorrelease
spec:
  gif:
    link: https://media.giphy.com/media/3PhKYCVdMi87u/giphy.gif
EOF"
wait
pe "$KUBECTL create -f releases/operatorrelease.yaml"
wait

# ----------------------------- Release
p "# Do release 2, with INVERTED business decisions"
pe "cat > releases/satisfied.yaml <<EOF
apiVersion: wigm.carson-anderson.com/v1
kind: WigmGif
metadata:
  name: satisfied
spec:
  gif:
    link: https://media.giphy.com/media/l2JegGMtnxw0Nq3pC/giphy.gif
  service:
    create_cloud_lb: true
  ingress:
    enabled: false
EOF"
wait
pe "$KUBECTL create -f releases/satisfied.yaml"
wait


# ----------------------------- Curl the ingress
p "# Curl the cluster ingress controller with our wigm hostname"
pe "NODEIP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' wigm-operator-sdk_node_1)"
pe "curl -s -H 'Host: operatorrelease.wigm.carson-anderson.com' $NODEIP"
wait

# ----------------------------- Proxy
p "# Proxy connection into a release 1 pod. Access at http://localhost:8888"
p "# Control-C to kill and continue"
pe "$KUBECTL port-forward svc/wigm-host-operatorrelease 8888:80"

# ----------------------------- Upgrade
p "# Upgrade release 1. Set a better name"
pe "cat > releases/operatorrelease.yaml <<EOF
apiVersion: wigm.carson-anderson.com/v1
kind: WigmGif
metadata:
  name: operatorrelease
spec:
  gif:
    name: Operators seek state!
    link: https://media.giphy.com/media/3PhKYCVdMi87u/giphy.gif
EOF"
wait
pe "$KUBECTL apply -f releases/operatorrelease.yaml"
wait

# ----------------------------- Curl the ingress
p "# Curl the cluster ingress controller with our wigm hostname"
pe "NODEIP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' wigm-operator-sdk_node_1)"
pe "curl -s -H 'Host: operatorrelease.wigm.carson-anderson.com' $NODEIP"
wait

# ----------------------------- Proxy
p "# Proxy connection into a release 1 pod. Access at http://localhost:8888"
p "# Control-C to kill and continue"
pe "$KUBECTL port-forward svc/wigm-host-operatorrelease 8888:80"

# ----------------------------- Curl the ingress
p "# Curl the cluster ingress controller with our release 2 wigm hostname"
pe "NODEIP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' wigm-operator-sdk_node_1)"
pe "curl -s -H 'Host: satisfied.wigm.carson-anderson.com' $NODEIP"
wait

# ----------------------------- Curl the ingress
p "# Patch the resouce to enable ingress"
pe "$KUBECTL patch wigmgif satisfied -p '{\"spec\":{\"ingress\":{\"enabled\":true}}}' --type=merge"
wait

# ----------------------------- Curl the ingress
p "# Curl the cluster ingress controller with our release 2 wigm hostname"
pe "NODEIP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' wigm-operator-sdk_node_1)"
pe "curl -s -H 'Host: satisfied.wigm.carson-anderson.com' $NODEIP"
wait

# ----------------------------- Break something
p "# Accidentally break something"
pe "$KUBECTL delete svc wigm-host-operatorrelease"
wait

# ----------------------------- Proxy
p "# Proxy connection into a release 1 pod, Control-C to kill and continue"
pe "$KUBECTL port-forward svc/wigm-host-operatorrelease 8888:80"
wait

# ----------------------------- cleanup
p "# Demo done! Don't forget to clean up!"
p "# Clean up all wigm releases"
pe "$KUBECTL delete wigmgif --all"
wait

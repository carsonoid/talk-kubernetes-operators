#!/bin/bash

. ../demo-magic.sh -n

export TYPE_SPEED=402
export DEMO_COMMENT_COLOR=$CYAN

export KUBECTL="kubectl --kubeconfig=kubeconfig.yaml"

# Create New Instance
p "# Create new instance with a Custom resource"
pe "cat > releases/operator1.yaml <<EOF
apiVersion: wigm.carson-anderson.com/v1
kind: WigmGif
metadata:
  name: operator1
spec:
  gif:
    link: https://media.giphy.com/media/l2JegGMtnxw0Nq3pC/giphy.gif
  service:
    create_cloud_lb: true
  ingress:
    enabled: false
EOF"
wait
pe "$KUBECTL apply -f releases/operator1.yaml"
wait

# ----------------------------- Proxy
p "# Proxy http://localhost:8888 into pod"
p "# Control-C to kill and continue"
pe "$KUBECTL port-forward svc/wigm-host-operator1 8888:80"

# ----------------------------- Update Title
p "# Update title"
pe "cat > releases/operator1.yaml <<EOF
apiVersion: wigm.carson-anderson.com/v1
kind: WigmGif
metadata:
  name: operator1
spec:
  gif:
    title: \"Operators are awesome!\"
    link: https://media.giphy.com/media/l2JegGMtnxw0Nq3pC/giphy.gif
  service:
    create_cloud_lb: true
  ingress:
    enabled: false
EOF"
wait
pe "$KUBECTL apply -f releases/operator1.yaml"
wait


# ----------------------------- Enable Ingress
p "# Enable ingress"
pe "$KUBECTL patch wigmgif operator1 -p '{\"spec\":{\"ingress\":{\"enabled\":true}}}' --type=merge"
wait

# ----------------------------- Proxy
p "# Check out changes!"
p "# Proxy http://localhost:8888 into pod"
p "# Control-C to kill and continue"
pe "$KUBECTL port-forward svc/wigm-host-operator1 8888:80"

# ----------------------------- Add new release
p "# Create another instance"
pe "cat > releases/operator2.yaml <<EOF
apiVersion: wigm.carson-anderson.com/v1
kind: WigmGif
metadata:
  name: operator2
spec:
  gif:
    title: \"Operators seek state!\"
    link: https://media.giphy.com/media/3PhKYCVdMi87u/giphy.gif
  service:
    create_cloud_lb: false
  ingress:
    enabled: false
EOF"
wait
pe "$KUBECTL apply -f releases/operator2.yaml"
wait

# ----------------------------- Report configuration
p "# Report configuration"
pe "$KUBECTL get wigmgif"
wait

# ----------------------------- Delete
p "# Delete"
pe "$KUBECTL delete wigmgif --all"
wait

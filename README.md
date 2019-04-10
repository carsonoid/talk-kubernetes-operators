# Kubernetes Operators

## Slides

Slides are found here in PDF format

## WIGM

WIGM is a sample project that is designed to help users understand why operators are useful. The project provides folders for 4 different deployment methods:

* `yaml`
* `helm`
* `metacontroller`
* `operator-sdk`

Each folder contains the following common files:

#### `start-and-watch.sh`

* Bootstraps a Kubernetes cluster with docker-compose using [k3s](https://github.com/rancher/k3s)
* Configures it for the demo
* Watches applicable resources

After the cluster has started. There will be a `kubeconfig.yaml` file in the current directory. `kubectl` can be used on the user's machine to access the demo cluster at any time by passing a flag (`kubectl --kubeconfig=kubeconfig.yaml`) or exporting an environment variable (`export KUBECONFIG=./kubeconfig.yaml`).

#### `demo.sh`

An automatic demo script which walks the user through a few WIGM deployments with the demo architecture

#### `cleanup.sh`

Brings down the docker resources and cleans up the docker volumes. NOTE: This deletes all cluster state

## Try It Out

Run a test cluster and demo for any deployment method

1. `cd METHODFOLDER`
2. `./start-and-watch.sh`
3. In a new terminal: `./demo.sh`
4. Optional: Run extra `kubectl` commands against the demo cluster if desired
5. Optional: `./cleanup.sh`

# Attach to node k3s

This is a tool which setups a shell environment with a SOCKS5 proxy to a node's k3s API server and a copy of the node's kubeconfig. This allows you to easily debug the node from the comfort of your own machine.

## Installation

You'll need to have the latest version of Go installed. Then you simply run:

```sh
go install github.com/waggle-sensor/attach-to-node-k3s@latest
```

You'll also need kubectl installed. (Either explicitly or as part of tools like Rancher Desktop.)

## Usage

You just run:

```sh
attach-to-node-k3s vsn
```

Now you should be in a new shell and can directly interact with the node:

```sh
# Should show you pods running on the *node*!
kubectl get pods
```

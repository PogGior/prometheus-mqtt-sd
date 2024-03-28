# Deploy examples

This directory contains examples of how to run the plugin on a container in a complex environment including Prometheus integration.

To integrate the plugin with Prometheus, remember to configure the plugin with the `output.file` argument and the Prometheus with the `file_sd_config` to read the file. 

It is important also to run the plugin container with user `nobody` to avoid permission issues. The user nobody is the default user for the Prometheus container and has the identifier `65534`.

The examples are:
- k8s: example of how to run the plugin in a Kubernetes environment. The installation includes a Mosquitto statefulset exposing a NodePort on 31883 and a Prometheus daemonset already configured with prometheus-mqtt-sd sidecar exposing a NodePort on 30090. The environment can be installed with kustomize tool with the following command:
```bash
kubectl apply -k k8s
```
- docker-compose: *COMING SOON*
 
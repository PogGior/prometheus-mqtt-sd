# Prometheus Mqtt SD

Prometheus MQTT Service Discovery is a simple plugin for Prometheus that uses MQTT to dynamically update the list of targets.

The implementation follows the guidelines specified posted by the Prometheus team [here](https://prometheus.io/blog/2018/07/05/implementing-custom-sd/) and it is based on the [file_sd_config](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#file_sd_config) supported Service Discovery mechanism.

## Use case

The main use case for prometheus-mqtt-sd is the IOT industry, where MQTT is one of the main protocols used for devices remote control.

This plugin allows you to dynamically update the list of targets in central or distributed monitoring system based on the devices that are connected to the MQTT broker.

It can be useful to add a target when a device is connected or its status changes. It can also be useful to enable monitoring depht, for example, when a device fail with error, the list of its targets can be updated to debug the problem.

## How it works

The plugin subscribes to a MQTT topic and listens for messages that contain the list of targets in JSON format.

When a message is received, the plugin writes the list of targets to a file in the format expected by Prometheus. The file is then read by Prometheus by the file_sd_config and the list of targets is updated.

The MQTT message format is the following:

```json
[
    {
        "targets": [],
        "labels": {},
        "source": ""
    }
]
```

- `targets` is a list of strings that represent the targets addresses to be monitored.
- `labels` is a dictionary of strings that represent the labels to be attached to the targets.
- `source` is a string that represents the source of the targets. Source is an unique identifier for a group of targets. To update the list of targets, the plugin uses the source to identify the targets that have to be updated. If the source is not present in the message, the plugin will not update the list of targets for that source.

An example of a message can be found [here](fixtures/example-message.json).

## Installation

The plugin can be builded using the following command:

```bash
make build
```

And runned locally with:

```bash
make run
```

To run the plugin in a container, you can use the following command:

```bash
make docker-image
```

Examples of how to run the plugin on a container can be found [here](deploy).

## Configuration



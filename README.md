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

And saved the image in a tar file with:

```bash
make docker-save
```

Examples of how to run the plugin on a container in a complex environment including Prometheus integration can be found [here](deploy).

## Configuration

The plugin can be first configures by providing two startup arguments:

```bash
./prometheus-mqtt-sd --config.file config.yaml --output.file output.json
```

- `output.file`: the file where the list of targets will be written in the format expected by Prometheus. The file is read by Prometheus by the file_sd_config.
- `config.file`: the file where the configuration of the plugin is stored. The configuration file is in YAML format and an example is the following:

```yaml
address: ssl://0.0.0.0:8883
topic: topic.test
client_id: prometheus-mqtt-sd
username: mqtt
password: mqtt
tls:
  ca_file: /etc/ssl/certs/ca.crt
  cert_file: /etc/ssl/certs/client.crt
  key_file: /etc/ssl/private/client.key
  insecure_skip_verify: false
```

- `address`: the address of the MQTT broker.
- `topic`: the topic where the plugin listens for messages.
- `client_id`: the client id used to connect to the MQTT broker.
- `username`: the username used to connect to the MQTT broker.
- `password`: the password used to connect to the MQTT broker.
- `tls`: the TLS configuration used to connect to the MQTT broker.
  - `ca_file`: the path to the CA file.
  - `cert_file`: the path to the client certificate file.
  - `key_file`: the path to the client key file.
  - `insecure_skip_verify`: a boolean that indicates if the TLS verification should be skipped.

Other examples of configuration files can be found [here](fixtures).

## Testing

The plugin can be tested using the following command:

```bash
make test
```




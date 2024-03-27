package util_test

import (
	"prometheus-mqtt-sd/cmd/prometheus-mqtt-sd/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testFileDir = "../../../fixtures"
)

func TestNewTlsConfig(t *testing.T) {
	// Test case 1: CAFile is not empty
	config := util.TLSConfig{
		CAFile:             testFileDir + "/ca.crt",
		InsecureSkipVerify: false,
		CertFile:           testFileDir + "/client.crt",
		KeyFile:            testFileDir + "/client.key",
	}
	tlsConfig, err := config.NewTlsConfig()
	assert.Nil(t, err)
	assert.NotNil(t, tlsConfig.RootCAs)

	// Test case 2: InsecureSkipVerify is true
	config = util.TLSConfig{
		CAFile:             "",
		InsecureSkipVerify: true,
		CertFile:           testFileDir + "/client.crt",
		KeyFile:            testFileDir + "/client.key",
	}
	tlsConfig, err = config.NewTlsConfig()
	assert.Nil(t, err)
	assert.True(t, tlsConfig.InsecureSkipVerify)

	// Test case 3: CertFile and KeyFile are not empty
	config = util.TLSConfig{
		CAFile:             "",
		InsecureSkipVerify: false,
		CertFile:           testFileDir + "/client.crt",
		KeyFile:            testFileDir + "/client.key",
	}
	tlsConfig, err = config.NewTlsConfig()
	assert.Nil(t, err)
	assert.NotNil(t, tlsConfig.Certificates)
	assert.Equal(t, 1, len(tlsConfig.Certificates))
}

func TestLoadConfig(t *testing.T) {

	// Test case 1: Valid config file 
	file := testFileDir + "/config-simple.yaml"
	expectedConfig := &util.SdConfig{
		Address:  "tcp://0.0.0.0:1883",
		Topic:    "topic.test",
		ClientID: "prometheus-mqtt-sd",
		Username: "",
		Password: "",
		TLSConfig: nil,
	}
	config, err := util.LoadConfig(file)
	assert.Nil(t, err)
	assert.Equal(t, expectedConfig, config)

	// Test case 2: Valid config file with TLS
	file = testFileDir + "/config-tls.yaml"
	expectedConfig = &util.SdConfig{
		Address:  "ssl://0.0.0.0:8883",
		Topic:    "topic.test",
		ClientID: "prometheus-mqtt-sd",
		Username: "mqtt",
		Password: "mqtt",
		TLSConfig: &util.TLSConfig{
			CAFile:             "/etc/ssl/certs/ca.crt",
			InsecureSkipVerify: false,
			CertFile:           "/etc/ssl/certs/client.crt",
			KeyFile:            "/etc/ssl/private/client.key",
		},
	}
	config, err = util.LoadConfig(file)
	assert.Nil(t, err)
	assert.Equal(t, expectedConfig, config)

	// Test case 3: Invalid config file
	file = testFileDir + "/invalid_config.yaml"
	_, err = util.LoadConfig(file)
	assert.NotNil(t, err)
}

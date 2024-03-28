package util

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kit/log"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type SdConfig struct {
	Address   string     `yaml:"address"`
	Topic     string     `yaml:"topic"`
	ClientID  string     `yaml:"client_id"`
	Username  string     `yaml:"username"`
	Password  string     `yaml:"password"`
	TLSConfig *TLSConfig `yaml:"tls"`
}

type TLSConfig struct {
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
	CAFile             string `yaml:"ca_file"`
	CertFile           string `yaml:"cert_file"`
	KeyFile            string `yaml:"key_file"`
}

func LoadConfig(file string) (*SdConfig, error) {
	var config SdConfig

	configFile, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (config *SdConfig) InitDiscovery(logger log.Logger) (*discovery, error) {

	client, err := config.initClient()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cd := &discovery{
		topic:  config.Topic,
		client: client,
		logger: logger,
	}

	return cd, nil
}

func (config *SdConfig) initClient() (mqtt.Client, error) {

	opts := mqtt.NewClientOptions().AddBroker(config.Address)
	opts.SetClientID(config.ClientID)
	if config.Username != "" {
		opts.SetUsername(config.Username)
	}
	if config.Password != "" {
		opts.SetPassword(config.Password)
	}
	if config.TLSConfig != nil {
		tlsConfig, err := config.TLSConfig.NewTlsConfig()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		opts.SetTLSConfig(tlsConfig)
	}

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}

func (config TLSConfig)NewTlsConfig() (*tls.Config, error) {
	tlsConfig := &tls.Config{}
	if config.CAFile != "" {
		certpool := x509.NewCertPool()

		ca, err := os.ReadFile(config.CAFile)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		certpool.AppendCertsFromPEM(ca)
		tlsConfig.RootCAs = certpool
	}
	if config.InsecureSkipVerify {
		tlsConfig.InsecureSkipVerify = true
	}
	if config.CertFile != "" && config.KeyFile != "" {
		clientKeyPair, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		tlsConfig.Certificates = []tls.Certificate{clientKeyPair}
	}
	return tlsConfig, nil
}

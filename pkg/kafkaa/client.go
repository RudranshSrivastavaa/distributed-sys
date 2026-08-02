package kafkaa

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/xdg-go/scram"
)

type Client struct {
	config Config
}

func NewClient(config Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) saramaConfig() *sarama.Config {

	cfg := sarama.NewConfig()

	cfg.ClientID = c.config.ClientID
	cfg.Version = c.config.Version

	cfg.Producer.RequiredAcks = c.config.Producer.RequiredAcks
	cfg.Producer.Retry.Max = c.config.Producer.RetryMax
	cfg.Producer.Return.Successes = c.config.Producer.ReturnSuccesses

	cfg.Consumer.Offsets.Initial = c.config.Consumer.InitialOffset

	if c.config.TLS.Enabled {

		caCert, err := os.ReadFile(c.config.TLS.CAFile)
		log.Println("Loading CA from:", c.config.TLS.CAFile)
		if err != nil {
			panic(err)
		}

		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(caCert) {
			panic("failed to append CA certificate")
		}

		cfg.Net.TLS.Enable = true
		cfg.Net.TLS.Config = &tls.Config{
			RootCAs:    pool,
			MinVersion: tls.VersionTLS12,
		}
	}

	if c.config.SASL.Enabled {

		cfg.Net.SASL.Enable = true

		cfg.Net.SASL.User = c.config.SASL.Username

		cfg.Net.SASL.Password = c.config.SASL.Password

		cfg.Net.SASL.Handshake = true

		cfg.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512

		cfg.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &XDGSCRAMClient{
				HashGeneratorFcn: scram.SHA512,
			}
		}
	}

	return cfg
}

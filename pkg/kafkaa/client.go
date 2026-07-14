package kafkaa

import "github.com/IBM/sarama"

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

	cfg.Producer.RequiredAcks =
		c.config.Producer.RequiredAcks

	cfg.Producer.Retry.Max =
		c.config.Producer.RetryMax

	cfg.Producer.Return.Successes =
		c.config.Producer.ReturnSuccesses

	cfg.Consumer.Offsets.Initial =
		c.config.Consumer.InitialOffset

	return cfg

}
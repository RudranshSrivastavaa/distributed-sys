package kafkaa

import (
	"context"

	"github.com/IBM/sarama"
)


type Consumer interface {

	Start(ctx context.Context) error

	Close() error
}

func (c *Client) NewConsumerGroup() (sarama.ConsumerGroup,error) {


	return sarama.NewConsumerGroup(

		c.config.Brokers,

		c.config.Consumer.GroupID,

		c.saramaConfig(),
	)

}
package kafkaa

import (
	"log"

	"github.com/IBM/sarama"
)

func (c *Client) NewProducer() (sarama.SyncProducer,error) {
	log.Printf("Kafka brokers in producer: %+v", c.config.Brokers)
	return sarama.NewSyncProducer(

		c.config.Brokers,

		c.saramaConfig(),
	)

}
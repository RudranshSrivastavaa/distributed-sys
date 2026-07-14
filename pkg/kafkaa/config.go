package kafkaa

import (
	"github.com/IBM/sarama"
)

type Config struct {
	Brokers []string

	ClientID string

	Version sarama.KafkaVersion

	Producer ProducerConfig

	Consumer ConsumerConfig
}

type ProducerConfig struct {
	RequiredAcks sarama.RequiredAcks

	RetryMax int

	ReturnSuccesses bool
}

type ConsumerConfig struct {
	GroupID string

	InitialOffset int64
}

func DefaultConfig() Config {

	return Config{

		Brokers: []string{
			"localhost:29092",
		},

		ClientID: "distributed-commerce",

		Version: sarama.V4_0_0_0,

		Producer: ProducerConfig{

			RequiredAcks: sarama.WaitForAll,

			RetryMax: 5,

			ReturnSuccesses: true,
		},

		Consumer: ConsumerConfig{

			GroupID: "distributed-commerce",

			InitialOffset: sarama.OffsetNewest,
		},
	}

}
package kafkaa

import (
	"github.com/IBM/sarama"
)

type TLSConfig struct {
	Enabled bool

	CAFile string
}

type SASLConfig struct {
	Enabled bool

	Username string

	Password string

	Mechanism sarama.SASLMechanism
}

type Config struct {
	Brokers []string

	ClientID string

	TLS TLSConfig

	SASL SASLConfig

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
			"localhost:9094",
		},

		ClientID: "distributed-commerce",

		TLS: TLSConfig{
			Enabled: false,
			CAFile:  "",
		},

		SASL: SASLConfig{

			Enabled: false,
		},

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

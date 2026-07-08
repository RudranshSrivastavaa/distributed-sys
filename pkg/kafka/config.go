package kafka

type Config struct {

	Brokers []string

	ClientID string
}

func DefaultConfig() Config {

	return Config{

		Brokers: []string{
			"localhost:9092",
		},

		ClientID: "distributed-commerce",
	}

}
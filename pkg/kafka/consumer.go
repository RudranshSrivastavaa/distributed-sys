package kafka

import "context"

type consumer struct {

	config Config

	topic Topic

	groupID string

	handler MessageHandler

	reader *kafkago.Reader
}

func NewConsumer(config Config,topic Topic,groupID string,handler MessageHandler) Consumer {

	reader := kafkago.NewReader(
		kafkago.ReaderConfig{
			Brokers: config.Brokers,
			GroupID: groupID,
			Topic: string(topic),
		},
	)
	return &consumer{

		config: config,

		topic: topic,

		groupID: groupID,

		handler: handler,

		reader: reader,
	}

}

func (c *consumer) Start(ctx context.Context) error {

	for {

		select {

		case <-ctx.Done():

			return nil

		default:

		}

		msg, err := c.reader.FetchMessage(
			ctx,
		)

		if err != nil {

			return err

		}

		message := Message{

			Topic: c.topic,

			Key: msg.Key,

			Value: msg.Value,

			Timestamp: msg.Time,

			Headers: headersToMap(
				msg.Headers,
			),
		}

		err = c.handler.Handle(
			ctx,
			message,
		)

		if err != nil {

			continue

		}

		if err := c.reader.CommitMessages(
			ctx,
			msg,
		); err != nil {

			return err

		}

	}

}

func (c *consumer) Close() error {

	return c.reader.Close()

}

func headersToMap(headers []kafkago.Header) map[string]string {

	result := make(map[string]string)

	for _, h := range headers {
		result[h.Key] = string(h.Value)
	}
	
	return result
}
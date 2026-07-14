package kafkaa

// import (
// 	"context"

// 	"github.com/IBM/sarama"
// )

// type consumer struct {

// 	group sarama.ConsumerGroup

// 	topic Topic

// 	handler EventHandler
// }

// func NewConsumer(

// 	client *Client,

// 	topic Topic,

// 	handler EventHandler,

// ) (Consumer, error) {

// 	group, err := client.NewConsumerGroup()

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &consumer{
// 		group: group,
// 		topic: topic,
// 		handler: handler,
// 	}, nil
// }

// func (c *consumer) Start(

// 	ctx context.Context,

// ) error {

// 	for {

// 		err := c.group.Consume(
// 			ctx,
// 			[]string{
// 				c.topic.Name,
// 			},
// 			&consumerHandler{
// 				ctx: ctx,
//           		handler: c.handler,
// 			},
// 		)
// 		if err != nil {
// 			return err
// 		}

// 		if ctx.Err() != nil {
// 			return nil
// 		}

// 	}

// }

// func (c *consumer) Close() error {

// 	return c.group.Close()

// }
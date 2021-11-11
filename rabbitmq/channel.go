package rabbitmq

import (
	"github.com/streadway/amqp"
)

type BasicQos struct {
	PrefetchCount int
	PrefetchSize  int
}

type Channel struct {
	*amqp.Channel
}

func NewChannel(client *Client) (*Channel, error) {
	channel, err := client.conn.Channel()
	if err != nil {
		return nil, err
	}
	ch := &Channel{
		Channel: channel,
	}
	return ch, nil
}

func (ch *Channel) qos(qos *BasicQos) error {
	return ch.Qos(qos.PrefetchCount, qos.PrefetchSize, false)
}

func (ch *Channel) privateQueueDeclare(queue *Queue) error {
	_, err := ch.QueueDeclare(
		queue.Name,
		queue.Durable,
		false,
		false,
		false,
		nil,
	)
	return err

}

package rabbitmq

import (
	"github.com/streadway/amqp"
	"sync"
)

type PublisherConfig struct {
	Exchange   Exchange `mapstructure:"exchange"`
	RoutingKey string   `mapstructure:"routing_key"`
	Queue      Queue    `mapstructure:"queue"`
}

type Publisher struct {
	client      *Client
	channel     *Channel
	Config      *PublisherConfig
	initialized bool
	m           sync.Mutex
}

func NewPublisher(cfg *PublisherConfig, client *Client) *Publisher {
	return &Publisher{Config: cfg, client: client}
}

func (p *Publisher) WithClient(client *Client) {
	p.client = client
}

func (p *Publisher) init() error {
	if p.client == nil {
		return clntOfPublisherCannotBeNil
	}
	p.m.Lock()
	defer p.m.Unlock()
	if p.initialized == true {
		return nil
	}
	channel, err := NewChannel(p.client)
	if err != nil {
		return err
	}
	if err := channel.ExchangeDeclare(
		p.Config.Exchange.Name,
		p.Config.Exchange.Type,
		p.Config.Exchange.Durable,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}
	if queue := p.Config.Queue; queue.Name != "" {
		if err := channel.privateQueueDeclare(&queue); err != nil {
			return err
		}
		if err := channel.QueueBind(queue.Name, p.Config.RoutingKey, p.Config.Exchange.Name, false, nil); err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	p.channel = channel
	p.initialized = true
	return nil
}

// recover 恢复 publisher 的 client 和 channel
func (p *Publisher) recover() error {
	if err := p.client.recover(); err != nil {
		return err
	}
	channel, err := NewChannel(p.client)
	if err != nil {
		return err
	}
	p.channel = channel
	return nil
}

func (p *Publisher) Publish(msg []byte) error {
	if p.client.closedVoluntarily {
		return clntClosed
	}
	if !p.initialized {
		if err := p.init(); err != nil {
			return err
		}
	}
	err := p.channel.Publish(
		p.Config.Exchange.Name,
		p.Config.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plan",
			Body:        msg,
		},
	)
	if err != nil {
		// 非连接错误直接返回
		if !isConnectionError(err) {
			return err
		}
		// 否则恢复连接
		if err := p.recover(); err != nil {
			return err
		}
		// 恢复连接后重新推送
		return p.Publish(msg)
	}
	return nil
}

package rabbitmq

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"sync"
)

type RabbitMQ struct {
	Options    *ClientOptions `mapstructure:"options"`
	cClient    *Client
	pClient    *Client
	Consumers  Consumers  `mapstructure:"consumers"`
	Publishers Publishers `mapstructure:"publishers"`
	m          sync.Mutex
}

func New(opts *ClientOptions, consumers Consumers, publishers Publishers) *RabbitMQ {
	return &RabbitMQ{
		Options:    opts,
		Consumers:  consumers,
		Publishers: publishers,
	}
}

func NewFromViper(v *viper.Viper) (*RabbitMQ, error) {
	return NewFromViperByKey(v, "rabbitmq")
}

func NewFromViperByKey(v *viper.Viper, key string) (*RabbitMQ, error) {
	mq := new(RabbitMQ)
	if err := v.UnmarshalKey(key, mq); err != nil {
		return nil, err
	}
	return mq, nil
}

// NewFromCfg 从已有配置生成 MQ 实例, 配置应当从 viper 中获取
//   e.g.
//   cfg := viper.GetStringMap("your-rabbitmq")
//   mq, err := NewFromCfg(cfg)
//func NewFromCfg(cfg map[string]interface{}) (*RabbitMQ, error) {
//	return parseConfig(cfg)
//}

func (mq *RabbitMQ) connectCClient() error {
	if mq.cClient != nil {
		return nil
	}
	c, err := NewClient(mq.Options)
	if err != nil {
		return err
	}
	mq.cClient = c
	return nil
}

func (mq *RabbitMQ) connectPClient() error {
	if mq.pClient != nil {
		return nil
	}
	c, err := NewClient(mq.Options)
	if err != nil {
		return err
	}
	mq.pClient = c
	return nil
}

func (mq *RabbitMQ) getConsumer(consumerName string) (*Consumer, error) {
	for _, consumer := range mq.Consumers {
		if consumer.Config.Name == consumerName {
			return consumer, nil
		}
	}
	return nil, fmt.Errorf("no consumer named %s", consumerName)
}

func (mq *RabbitMQ) getPublisher(exchange, routingKey string) (*Publisher, error) {
	for _, publisher := range mq.Publishers {
		if publisher.Config.Exchange.Name == exchange && publisher.Config.RoutingKey == routingKey {
			return publisher, nil
		}
	}
	return nil, fmt.Errorf("no publisher binding exchange: %s and routing_key: %s", exchange, routingKey)
}

func (mq *RabbitMQ) Consume(consumerName string, callback func(dely *amqp.Delivery)) error {
	if err := mq.connectCClient(); err != nil {
		return err
	}
	consumer, err := mq.getConsumer(consumerName)
	if err != nil {
		return err
	}
	consumer.withClient(mq.cClient)
	return consumer.Consume(callback)
}

func (mq *RabbitMQ) Publish(exchange, routingKey string, msg []byte) error {
	if err := mq.connectPClient(); err != nil {
		return err
	}
	publisher, err := mq.getPublisher(exchange, routingKey)
	if err != nil {
		return err
	}
	publisher.WithClient(mq.pClient)
	return publisher.Publish(msg)
}

func (mq *RabbitMQ) Close() error {
	if mq.cClient != nil {
		return mq.cClient.Close()
	}
	if mq.pClient != nil {
		return mq.pClient.Close()
	}
	return nil
}

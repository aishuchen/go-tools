package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ConsumerConfig struct {
	Name          string `mapstructure:"name"`
	Queue         Queue  `mapstructure:"queue"`
	AutoAck       bool   `mapstructure:"auto_ack"`
	PrefetchCount int    `mapstructure:"prefetch_count"` // 注意 prefetch_count 和 prefetch_size 在 auto_ack 为 true 的情况下会失效
	PrefetchSize  int    `mapstructure:"prefetch_size"`
}

type Consumer struct {
	client      *Client
	channel     *Channel
	Config      *ConsumerConfig
	initialized bool
	m           sync.Mutex
}

func NewConsumer(cfg *ConsumerConfig, client *Client) *Consumer {
	return &Consumer{Config: cfg, client: client}
}

func (c *Consumer) withClient(client *Client) {
	c.client = client
}

func (c *Consumer) init() error {
	if c.client == nil {
		return clntOfConsumerCannotBeNil
	}
	c.m.Lock()
	defer c.m.Unlock()
	if c.initialized {
		return nil
	}
	channel, err := NewChannel(c.client)
	if err != nil {
		return err
	}
	if c.Config.AutoAck {
		return nil
	}
	prefetchCount := c.Config.PrefetchCount
	prefetchSize := c.Config.PrefetchSize
	if prefetchCount == 0 {
		prefetchCount = 1
	}
	qosOpts := &BasicQos{
		PrefetchCount: prefetchCount,
		PrefetchSize:  prefetchSize,
	}
	if err := channel.qos(qosOpts); err != nil {
		return err
	}
	queue := c.Config.Queue
	if err := channel.privateQueueDeclare(&queue); err != nil {
		return err
	}
	c.channel = channel
	c.initialized = true
	return nil
}

func (c *Consumer) recover() error {
	if err := c.client.recover(); err != nil {
		return err
	}
	channel, err := NewChannel(c.client)
	if err != nil {
		return err
	}
	c.channel = channel
	return nil
}

func (c *Consumer) Consume(callback func(dely *amqp.Delivery)) error {
	if c.client.closedVoluntarily {
		return clntClosed
	}
	if !c.initialized {
		if err := c.init(); err != nil {
			return err
		}
	}
	delyChan, err := c.channel.Consume(c.Config.Queue.Name, c.Config.Name, c.Config.AutoAck, false, false, false, nil)
	if err != nil {
		return err
	}
	sigCh := make(chan os.Signal)
	notifyClose := make(chan *amqp.Error)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL)
	logger.Info(fmt.Sprintf("%s consuming...", c.Config.Name))
	defer logger.Sync()
	c.channel.NotifyClose(notifyClose)
	for {
		select {
		case dely := <-delyChan:
			callback(&dely)
		case sig := <-sigCh:
			if err := c.channel.Close(); err != nil {
				return err
			}
			return ConsumerExitedAbnormally(fmt.Sprintf("consumer exit by signal: %v", sig))
		case err := <-notifyClose:
			logger.Error(fmt.Sprintf("%s consume failed, Reason: %v", c.Config.Name, err))
			if !isConnectionError(err) {
				return err
			}
			if err := c.recover(); err != nil {
				return err
			}
			return c.Consume(callback)
		}
	}
}

package rabbitmq

import "errors"

type Queue struct {
	Name    string `mapstructure:"name"`
	Durable bool   `mapstructure:"durable"`
}

type Exchange struct {
	Name    string `mapstructure:"name"`
	Type    string `mapstructure:"type"`
	Durable bool   `mapstructure:"durable"`
}

type ConsumerExitedAbnormally string

func (e ConsumerExitedAbnormally) Error() string {
	return string(e)
}

type ConsumerName string
type PublisherName string

type Consumers []*Consumer
type Publishers []*Publisher

var clntOfPublisherCannotBeNil = errors.New("client of Publisher cannot be nil, call WithClient() first")
var clntOfConsumerCannotBeNil = errors.New("client of Consumer cannot be nil, call WithClient() first")
var clntClosed = errors.New("client already closed")

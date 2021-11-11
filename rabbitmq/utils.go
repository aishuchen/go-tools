package rabbitmq

import "github.com/streadway/amqp"

func isConnectionError(err error) bool {
	if serr, ok := err.(*amqp.Error); ok {
		return serr.Code == amqp.ConnectionForced || serr.Code == amqp.ChannelError
	}
	return false
}

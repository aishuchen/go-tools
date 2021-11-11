package rabbitmq

var testOptions = &ClientOptions{
	DSN:            "amqp://guest:guest@10.0.6.37:5672//",
	ReconnectDelay: 3,
}

var testIncorrectOptions = &ClientOptions{
	DSN:            "amqp://guest:guest@no-exists-host:5672//",
	ReconnectDelay: 3,
}

var testQueue = Queue{
	Name:    "my-queue",
	Durable: true,
}

var testConsumerConfig = &ConsumerConfig{
	Name:          "my-consumer",
	Queue:         testQueue,
	AutoAck:       false,
	PrefetchSize:  0,
	PrefetchCount: 1,
}

var testExchange = Exchange{
	Name:    "my-exchange",
	Type:    "direct",
	Durable: true,
}

var testRoutingKey = "my-routingkey"

var testPublisherConfig = &PublisherConfig{
	Exchange:   testExchange,
	RoutingKey: testRoutingKey,
	Queue:      testQueue,
}

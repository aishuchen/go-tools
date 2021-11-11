package rabbitmq

import (
	"github.com/streadway/amqp"
	"sync"
	"testing"
)

func TestConsumer_Consume(t *testing.T) {
	clnt, err := NewClient(testOptions)
	if err != nil {
		t.Fatal(err)
	}
	consumer := NewConsumer(testConsumerConfig, clnt)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err = consumer.Consume(func(dely *amqp.Delivery) {
			t.Log(string(dely.Body))
			if err := dely.Ack(false); err != nil {
				t.Log(err)
			}
		})
		wg.Done()
	}()
	wg.Wait()
	if err != nil {
		serr, ok := err.(ConsumerExitedAbnormally)
		if ok {
			t.Log(serr)
			return
		}
		t.Fatal(err)
	}
	if err := clnt.Close(); err != nil {
		t.Fatal(err)
	}
}

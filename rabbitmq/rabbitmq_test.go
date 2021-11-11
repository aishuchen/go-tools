package rabbitmq

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/aishuchen/go-tools/config"
	"github.com/aishuchen/go-tools/internal"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

func TestNewFromViper(t *testing.T) {
	configFilePath := internal.GetTestConfigFile()
	if err := config.SetGlobalConfig(configFilePath); err != nil {
		t.Fatal(err)
	}
	mq, err := NewFromViper(viper.GetViper())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(mq)
	mq.Close()
}

func testNew() *RabbitMQ {
	consumers := make(Consumers, 1)
	consumers[0] = NewConsumer(testConsumerConfig, nil)
	publishers := make(Publishers, 1)
	publishers[0] = NewPublisher(testPublisherConfig, nil)
	mq := New(testOptions, consumers, publishers)
	return mq
}

func TestNew(t *testing.T) {
	testNew()
}

func callback(dely *amqp.Delivery) {
	fmt.Println(string(dely.Body))
	dely.Ack(false)
}

// 注意：执行此测试用例会被阻塞
func TestRabbitMQ_Consume(t *testing.T) {
	mq := testNew()
	if err := mq.Consume("my-consumer", callback); err != nil {
		t.Logf("consume error: %v", err)
	}
	if err := mq.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestRabbitMQ_Publish(t *testing.T) {
	mq := testNew()
	defer mq.Close()
	if err := mq.Publish("my-exchange", "my-routingkey", []byte("hello world!")); err != nil {
		t.Fatal(err)
	}
	// 服务端手动重启，测试重连
	time.Sleep(10 * time.Second)
	if err := mq.Publish("my-exchange", "my-routingkey", []byte("hello world!")); err != nil {
		t.Fatal(err)
	}
}

func TestRabbitMQ_Publish_Consume(t *testing.T) {
	mq := testNew()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if err := mq.Publish("my-exchange", "my-routingkey", []byte("hello world!")); err != nil {
			t.Log(err)
		}
	}()
	go func() {
		if err := mq.Consume("my-consumer", callback); err != nil {
			t.Log(err)
			wg.Done()
		}
	}()
	wg.Wait()
	if err := mq.Close(); err != nil {
		t.Fatal(err)
	}
}

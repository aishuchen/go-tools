package rabbitmq

import (
	"testing"
)

func TestPublisher_Publish(t *testing.T) {
	clnt, err := NewClient(testOptions)
	if err != nil {
		t.Fatal(err)
	}
	pub := NewPublisher(testPublisherConfig, clnt)
	if err := pub.Publish([]byte("hello world")); err != nil {
		t.Fatal(err)
	}
	if err := clnt.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestPublisher_Publish2(t *testing.T) {
	clnt, err := NewClient(testOptions)
	if err != nil {
		t.Fatal(err)
	}
	pub := NewPublisher(testPublisherConfig, clnt)
	if err := clnt.Close(); err != nil {
		t.Fatal(err)
	}
	err = pub.Publish([]byte("hello world"))
	t.Log(err)
	if err != clntClosed {
		t.Fatal(err)
	}

}

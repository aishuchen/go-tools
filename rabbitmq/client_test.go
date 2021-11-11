package rabbitmq

import "testing"

func TestNewClientSuccess(t *testing.T) {
	clnt, err := NewClient(testOptions)
	if err != nil {
		t.Fatal(err)
	}
	clnt.Close()
}

func TestNewClientError(t *testing.T) {
	_, err := NewClient(testIncorrectOptions)
	if err == nil {
		t.Fatal("err should not be nil")
	}
}

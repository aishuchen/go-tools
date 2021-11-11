package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.hypers.com/server-go/tools/logging"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var logger = logging.DefaultLogger

const defaultReconnectDelay = 5

type ClientOptions struct {
	DSN            string `mapstructure:"dsn"`
	ReconnectDelay int    `mapstructure:"reconnect_delay"`
}

func (opts *ClientOptions) init() {
	if opts.ReconnectDelay == 0 {
		opts.ReconnectDelay = defaultReconnectDelay
	}
}

type Client struct {
	conn              *amqp.Connection
	opts              *ClientOptions
	closedVoluntarily bool // 主动关闭
	m                 sync.Mutex
}

func NewClient(opts *ClientOptions) (*Client, error) {
	clnt := new(Client)
	conn, err := connect(opts)
	if err != nil {
		return nil, err
	}
	clnt.conn = conn
	clnt.opts = opts
	return clnt, nil

}

func connect(opts *ClientOptions) (*amqp.Connection, error) {
	conn, err := amqp.Dial(opts.DSN)
	return conn, err
}

func (c *Client) Close() error {
	if c.closedVoluntarily { // 已关闭的 client 再次关闭不报错
		return nil
	}
	c.closedVoluntarily = true
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) IsClosed() bool {
	if c.conn != nil {
		return c.conn.IsClosed()
	}
	return true
}

// recover 恢复 client 与 server 的链接
func (c *Client) recover() error {
	if c.closedVoluntarily {
		return clntClosed
	}
	c.m.Lock()
	defer c.m.Unlock()
	if !c.IsClosed() { // 连接有效的情况下不做重连
		return nil
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL)
	connChan := make(chan *amqp.Connection)
	errChan := make(chan error)
	delay := time.Duration(c.opts.ReconnectDelay) * time.Second
	go func(connch chan<- *amqp.Connection, errch chan<- error) {
		for {
			select {
			case <-time.After(delay):
				logger.Info(fmt.Sprintf("reconnect after %v", delay))
				conn, err := connect(c.opts)
				if err == nil {
					connChan <- conn
					errChan <- err
					return
				}
			case sig := <-sigChan:
				connChan <- nil
				errChan <- fmt.Errorf("stop connection, exit by signal: %v", sig)
				return
			}
		}
	}(connChan, errChan)
	conn := <-connChan
	err := <-errChan

	close(connChan)
	close(errChan)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

module github.com/aishuchen/go-tools/rabbitmq

go 1.16

replace (
	github.com/aishuchen/go-tools/config => ../config
	github.com/aishuchen/go-tools/internal => ../internal
	github.com/aishuchen/go-tools/logging => ../logging
)

require (
	github.com/spf13/viper v1.8.1
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.7.0
	github.com/aishuchen/go-tools/config v0.0.4
	github.com/aishuchen/go-tools/internal v0.0.4
	github.com/aishuchen/go-tools/logging v0.0.1
)

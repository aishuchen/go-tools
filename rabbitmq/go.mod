module gitlab.hypers.com/server-go/tools/rabbitmq

go 1.16

replace (
	gitlab.hypers.com/server-go/tools/config => ../config
	gitlab.hypers.com/server-go/tools/internal => ../internal
	gitlab.hypers.com/server-go/tools/logging => ../logging
)

require (
	github.com/spf13/viper v1.8.1
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.7.0
	gitlab.hypers.com/server-go/tools/config v0.0.4
	gitlab.hypers.com/server-go/tools/internal v0.0.4
	gitlab.hypers.com/server-go/tools/logging v0.0.1
)

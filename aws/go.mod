module github.com/aishuchen/go-tools/aws

go 1.16

replace (
	github.com/aishuchen/go-tools/config => ../config
	github.com/aishuchen/go-tools/internal => ../internal
	github.com/aishuchen/go-tools/logging => ../logging

)

require (
	github.com/aws/aws-sdk-go-v2 v1.9.1
	github.com/aws/aws-sdk-go-v2/credentials v1.4.2
	github.com/aws/aws-sdk-go-v2/service/s3 v1.16.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.7.1
	github.com/spf13/viper v1.9.0
	github.com/aishuchen/go-tools/config v0.0.4
	github.com/aishuchen/go-tools/internal v0.0.4
	github.com/aishuchen/go-tools/logging v0.0.1
)

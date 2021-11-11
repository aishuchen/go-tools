package common

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/viper"
	"gitlab.hypers.com/server-go/tools/logging"
)

var logger = logging.DefaultLogger

const (
	AWS_ROLE_ARN                = "AWS_ROLE_ARN"
	AWS_WEB_IDENTITY_TOKEN_FILE = "AWS_WEB_IDENTITY_TOKEN_FILE"
)

// SimpleAWSConfig 是 aws.Config 的缩减版, 为了能从 viper 直接 unmarshal 到结构体
type SimpleAWSConfig struct {
	Region    string `mapstructure:"region"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

// NewAWSCfgFromViper 从viper生成 aws 配置, 返回 aws.Config 实例
func NewAWSCfgFromViper(v *viper.Viper, fns ...func(cfg *aws.Config)) (*aws.Config, error) {
	cfg := new(SimpleAWSConfig)
	if err := v.UnmarshalKey("aws", cfg); err != nil {
		return nil, err
	}
	return NewAWSCfg(cfg)
}

func NewAWSCfg(cfg *SimpleAWSConfig, options ...func(cfg *aws.Config)) (*aws.Config, error) {
	credsProvier := resolveCredsProvider(cfg)
	awsConfig := &aws.Config{
		Region:      cfg.Region,
		HTTPClient:  &http.Client{},
		Credentials: credsProvier,
	}
	for _, option := range options {
		option(awsConfig)
	}
	return awsConfig, nil
}

// resolveCredsProvider 从已有配置生成 CredentialsProvider,
// 有 access_key 和 secret_key (不为空字符串) 时, 返回 StaticCredentialsProvider;
// 环境变量中存在 AWS_ROLE_ARN 和 AWS_WEB_IDENTITY_TOKEN_FILE 时, 返回 stscreds.WebIdentityRoleProvider;
// 否则返回 ec2rolecreds.Provider
func resolveCredsProvider(cfg *SimpleAWSConfig) aws.CredentialsProvider {
	defer logger.Sync()
	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		logger.Info("New credentials provider with static.")
		return credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")
	}
	roleArn := os.Getenv(AWS_ROLE_ARN)
	tokenFile := os.Getenv(AWS_WEB_IDENTITY_TOKEN_FILE)
	if roleArn != "" && tokenFile != "" {
		tokenRetriever := stscreds.IdentityTokenFile(tokenFile)
		opts := sts.Options{Region: cfg.Region}
		client := sts.New(opts)
		logger.Info("New credentials provider with web identity role. ARN: " + roleArn)
		return stscreds.NewWebIdentityRoleProvider(client, roleArn, tokenRetriever)
	}
	logger.Info("New credentials provider with ec2 role")
	return ec2rolecreds.New()
}

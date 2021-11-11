package common

import (
	"testing"

	"github.com/aishuchen/go-tools/config"
	"github.com/aishuchen/go-tools/internal"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/spf13/viper"
)

func TestGetAWSCfgFromFile(t *testing.T) {
	configFilePath := internal.GetTestConfigFile()
	if err := config.SetGlobalConfig(configFilePath); err != nil {
		t.Fatal(err)
	}
	cfg, err := NewAWSCfgFromViper(viper.GetViper())
	if err != nil {
		t.Fatal(err)
	}
	switch typ := cfg.Credentials.(type) {
	case credentials.StaticCredentialsProvider:
		// pass
		t.Log(typ)
	default:
		t.Fatalf("CredentialsProvider should not be %v", typ)
	}
}

package common

import (
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/spf13/viper"
	"gitlab.hypers.com/server-go/tools/config"
	"gitlab.hypers.com/server-go/tools/internal"
	"testing"
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
	switch typ :=cfg.Credentials.(type) {
	case credentials.StaticCredentialsProvider:
		//pass
		t.Log(typ)
	default:
		t.Fatalf("CredentialsProvider should not be %v", typ)
	}
}
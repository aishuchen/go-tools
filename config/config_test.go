package config

import (
	"os"
	"testing"

	"github.com/aishuchen/go-tools/internal"
	"github.com/spf13/viper"
)

var configFilePath = internal.GetTestConfigFile()

func TestSetGlobalConfig(t *testing.T) {
	if err := SetGlobalConfig(configFilePath); err != nil {
		t.Fatal(err)
	}
	testVal := viper.GetString("env")
	if testVal != "test" {
		t.Fatalf(`val should be "test", but got "%s"`, testVal)
	}
}

func TestSetLocalConfig(t *testing.T) {
	v, err := SetLocalConfig(configFilePath)
	if err != nil {
		t.Fatal(err)
	}
	testVal := v.GetString("env")
	if testVal != "test" {
		t.Fatalf(`val should be "test", but got "%s"`, testVal)
	}
}

func TestReadInGlobalConfig(t *testing.T) {
	f, err := os.Open(configFilePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := ReadInGlobalConfig(f); err != nil {
		t.Fatal(err)
	}
	testVal := viper.GetString("env")
	if testVal != "test" {
		t.Fatalf(`val should be "test", but got "%s"`, testVal)
	}
}

func TestReadInLocalConfig(t *testing.T) {
	f, err := os.Open(configFilePath)
	if err != nil {
		t.Fatal(err)
	}
	v, err := ReadInLocalConfig(f)
	if err != nil {
		t.Fatal(err)
	}
	testVal := v.GetString("env")
	if testVal != "test" {
		t.Fatalf(`val should be "test", but got "%s"`, testVal)
	}
}

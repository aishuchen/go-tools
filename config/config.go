package config

import (
	"errors"
	"github.com/spf13/viper"
	"io"
	"os"
)

var CannotGetCfg = errors.New("cannot get config from viper")

func setConfig(v *viper.Viper, r io.Reader) error {
	v.SetConfigType("toml")
	if err := v.ReadConfig(r); err != nil {
		return err
	}
	return nil
}

// SetGlobalConfig 从配置文件配置全局 viper
func SetGlobalConfig(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	return setConfig(viper.GetViper(), f)
}

func ReadInGlobalConfig(r io.Reader) error  {
	return setConfig(viper.GetViper(), r)
}

// SetLocalConfig 从配置文件配置局部 viper, 与全局的 viper 隔离, 返回一个新的 Viper 实例和 error
func SetLocalConfig(filepath string) (*viper.Viper, error) {
	v := viper.New()
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	if err := setConfig(v, f); err != nil {
		return nil, err
	}
	return v, nil
}

func ReadInLocalConfig(r io.Reader) (*viper.Viper, error) {
	v := viper.New()
	if err := setConfig(v, r); err != nil {
		return nil, err
	}
	return v, nil
}
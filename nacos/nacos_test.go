package nacos

import (
	"fmt"
	"testing"
)

func TestNacos(t *testing.T) {
	client, _ := NewNacos("server-go", "server-go", "https://nacos-prd.hypers.cc", "", "DEFAULT_GROUP")
	content, _ := client.Get("mysql")
	fmt.Println(string(content))
}

type config struct {
	Mysql mysql
	App   app
}

type mysql struct {
	Host string
}
type app struct {
	Mode string
}

func TestNacosLoader(t *testing.T) {
	loader := NacosLoader{}
	c := new(config)
	loader.Load(c)
}

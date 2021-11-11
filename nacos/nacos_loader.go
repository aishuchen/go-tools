package nacos

import (
	"errors"
	"strings"

	"github.com/fatih/structs"
	"gopkg.in/yaml.v2"
)

// NacosLoader load config forom nacos
type NacosLoader struct {
	na *Nacos
}

// Load 从Nacos获取配置文件
// Naocs中的配置为yaml格式
// dataId为小写
// Example:
/*
type Config struct{
	App App
	Mysql Mysql
}
type App struct{
	Name string
}
type Mysql struct{
	Host string
}
则在nacos中应该有两个dataId
dataId App
Name: myapp

dataId Mysql
Host: 127.0.0.1

*/
func (n *NacosLoader) Load(s interface{}) error {
	if n.na == nil {
		return errors.New("nacos client is empty")
	}
	for _, field := range structs.Fields(s) {
		m := make(map[string]interface{})
		name := strings.ToLower(field.Name())
		data, err := n.na.Get(name)
		if string(data) == "" || err != nil {
			continue
		}
		if err := yaml.Unmarshal(data, &m); err != nil {
			field.Set(m)
		}
	}
	return nil
}

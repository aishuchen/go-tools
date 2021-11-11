package nacos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

const (
	token   = "/nacos/v1/auth/login"
	configs = "/nacos/v1/cs/configs"
)

type Nacos struct {
	Username  string
	Password  string
	Namespace string
	// https://nacos-prd.hypers.cc
	Endpoint string
	Group    string
	Client   *http.Client
	token    string
}

type tokenResp struct {
	AccessToken string `json:"accessToken"`
	TokenTtl    int    `json:"tokenTtl"`
	GlobalAdmin bool   `json:"globalAdmin"`
}

func (n *Nacos) getToken() (string, error) {
	var p tokenResp
	url := fmt.Sprintf("%s%s?username=%s&password=%s", n.Endpoint, token, n.Username, n.Password)
	body, err := n.do(url, "POST")
	if err != nil {
		return "", err
	}
	if err = json.Unmarshal(body, &p); err != nil {
		return "", fmt.Errorf("获取token失败")
	} else {
		return p.AccessToken, nil
	}
}

func NewNacos(username string, password string, endpoint string, namespace string, group string) (*Nacos, error) {
	n := &Nacos{
		Username:  username,
		Password:  password,
		Endpoint:  endpoint,
		Namespace: namespace,
		Group:     group,
		Client:    &http.Client{},
	}
	token, err := n.getToken()
	n.token = token
	return n, err
}

func NewNacosFromViper() (*Nacos, error) {
	n := &Nacos{
		Username:  viper.GetString("nacos.user"),
		Password:  viper.GetString("nacos.pass"),
		Endpoint:  viper.GetString("nacos.endpoint"),
		Namespace: viper.GetString("nacos.namespace"),
		Group:     viper.GetString("nacos.group"),
		Client:    &http.Client{},
	}
	token, err := n.getToken()
	n.token = token
	return n, err
}

func (n *Nacos) do(url string, method string) ([]byte, error) {
	req, _ := http.NewRequest(method, url, nil)
	resp, err := n.Client.Do(req)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	return content, err
}

func (n *Nacos) Get(dataId string) ([]byte, error) {
	url := fmt.Sprintf("%s%s?accessToken=%s&tenant=%s&group=%s&dataId=%s",
		n.Endpoint, configs, n.token, n.Namespace, n.Group, dataId)
	return n.do(url, "GET")
}

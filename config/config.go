package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var (
	C      *Config
	Prompt string
)

type Config struct {
	*sync.RWMutex
	Gpt struct {
		TextConfig  AuthConfig `json:"text_config"`
		ImageConfig AuthConfig `json:"image_config"`
	} `json:"gpt"`
	ContextStatus  bool   `json:"context_status"`
	BaseModel      string `json:"base_model"`
	KeepaliveRobot string `json:"keepalive_robot"`
}

func (c *Config) GetBaseModel() string {
	c.RLock()
	defer c.RUnlock()
	return c.BaseModel
}

func (c *Config) SetBaseModel(model string) {
	c.Lock()
	defer c.Unlock()
	c.BaseModel = model
}

type AuthConfig struct {
	ProxyUrl      string `json:"proxy_url"` //代理地址，不填使用官方地址
	AuthToken     string `json:"auth_token"`
	TriggerPrefix string `json:"trigger_prefix"`
}

func (c *Config) IsValid() bool {

	authConfigs := []AuthConfig{
		c.Gpt.TextConfig,
		c.Gpt.ImageConfig,
	}

	for _, authConfig := range authConfigs {
		if authConfig.AuthToken == "" || authConfig.TriggerPrefix == "" {
			return false
		}
	}
	return true
}

func InitConfig() {
	// 1. 读取 `config.json`
	data, err := os.ReadFile("./config/config.json")
	if err != nil {
		log.Fatalf("读取配置文件失败，请检查配置文件 `/config/config.json` 的配置, 错误信息: %+v\n", err)
	}
	config := Config{
		RWMutex: new(sync.RWMutex),
	}
	if err = json.Unmarshal(data, &config); err != nil {
		log.Fatalf("读取配置文件失败，请检查配置文件 `config.json` 的格式, 错误信息: %+v\n", err)
	}
	if !config.IsValid() {
		log.Fatal("配置文件校验失败，请检查 `config.json`")
	}
	C = &config
	// 2. 读取 prompt.txt
	prompt, err := os.ReadFile("./config/prompt.conf")
	if err != nil {
		log.Fatalf("读取配置文件失败，请检查配置文件 `prompt.conf` 的配置, 错误信息: %+v\n", err)
	}
	Prompt = string(prompt)
}

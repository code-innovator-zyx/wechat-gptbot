package config

import (
	"encoding/json"
	"log"
	"os"
)

var (
	C      *Config
	Prompt string
)

type Config struct {
	Gpt struct {
		TextConfig  AuthConfig `json:"text_config"`
		ImageConfig AuthConfig `json:"image_config"`
	} `json:"gpt"`
	ContextStatus bool   `json:"context_status"`
	BaseModel     string `json:"base_model"`
}

type AuthConfig struct {
	BaseURL       string `json:"base_url"`
	AuthToken     string `json:"auth_token"`
	TriggerPrefix string `json:"trigger_prefix"`
}

func (c *Config) IsValid() bool {

	authConfigs := []AuthConfig{
		c.Gpt.TextConfig,
		c.Gpt.ImageConfig,
	}

	for _, authConfig := range authConfigs {
		if authConfig.BaseURL == "" || authConfig.AuthToken == "" || authConfig.TriggerPrefix == "" {
			return false
		}
	}
	return true
}

func InitConfig() {
	// 1. 读取 `config.json`
	data, err := os.ReadFile("./config/config.json")
	if err != nil {
		log.Fatalf("读取配置文件失败，请检查配置文件 `config.json` 的配置, 错误信息: %+v\n", err)
	}
	config := Config{}
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
		log.Fatalf("读取配置文件失败，请检查配置文件 `prompt.txt` 的配置, 错误信息: %+v\n", err)
	}
	Prompt = string(prompt)
}

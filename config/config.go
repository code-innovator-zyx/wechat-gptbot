package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	C      *Config
	Prompt string
)

const (
	defaultNewCron     = "0 30 7 1/1 * ?"
	defaultWeatherCron = "0 0 8 1/1 * ?"
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
	CronConfig     struct {
		WeatherConfig WeatherCronConfig `json:"weather_config"`
		NewsConfig    NewsCronConfig    `json:"news_config"`
	}
}

// WeatherCronConfig 天气预报定时任务配置
type WeatherCronConfig struct {
	Users []struct {
		Name string `json:"name"` // 用户名
		City string `json:"city"` // 城市
	} `json:"users"`
	Spec string `json:"spec"` // cron 表达式
}

type NewsCronConfig struct {
	Users  []string // 用户名
	Groups []string // 群名称
	Spec   string   `json:"spec"` // cron 表达式
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

func (c *Config) CheckCronValid() {
	if c.CronConfig.WeatherConfig.Spec == "" {
		c.CronConfig.WeatherConfig.Spec = defaultWeatherCron
	}

	if c.CronConfig.NewsConfig.Spec == "" {
		c.CronConfig.NewsConfig.Spec = defaultNewCron
	}
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

	// 读取定时任务配置表
	cronConfig, err := os.ReadFile("./config/cron.json")
	if err != nil {
		log.Fatalf("读取配置文件失败，请检查配置文件 `cron.json` 的配置, 错误信息: %+v\n", err)
	}
	err = json.Unmarshal(cronConfig, &C.CronConfig)
	if err != nil {
		fmt.Println(err)
	}
	C.CheckCronValid()
}

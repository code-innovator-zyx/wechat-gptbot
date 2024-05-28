package core

import (
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"wechat-gptbot/config"
	"wechat-gptbot/core/gpt"
	"wechat-gptbot/core/handler"
	"wechat-gptbot/core/svc"
	"wechat-gptbot/logger"
	"wechat-gptbot/streamlit_web"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/4/8 17:12
* @Package:
 */
func Initialize() {
	// 初始化日志
	logger.InitLogrus(logger.Config{
		Level:      logrus.DebugLevel,
		ObjectName: "wechat-gptbot",
		WriteFile:  false,
	})
	// 初始化配置文件
	config.InitConfig()
	// 初始化会话上下文管理器
	handler.Context = svc.NewServiceContext()
	// 初始化提示词
	gpt.PromptMessage = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: config.Prompt,
	}
	go streamlit_web.RunStreamlit()
}

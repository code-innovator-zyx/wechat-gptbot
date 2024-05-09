package gpt

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"testing"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/5/9 10:34
* @Package:
 */
var clients *openAiClient

func init() {
	clients = &openAiClient{}
	clientConfig := openai.DefaultConfig("sk-yTGBVN2WlsMja5ADC879Fa6e1e044b22B07195EfC1A06dC4")
	clientConfig.BaseURL = "https://api.gpt.ge/v1"
	client := openai.NewClientWithConfig(clientConfig)
	clients.cs = map[string]*openai.Client{
		openai.GPT3Dot5Turbo: client,
	}
}
func Test_Chat(t *testing.T) {
	msgs := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "你是一个家庭厨师",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "今天晚上吃什么饭",
		},
	}
	t.Log(clients.createChat(context.Background(), openai.GPT3Dot5Turbo, msgs))
}

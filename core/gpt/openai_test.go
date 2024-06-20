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
		{Role: openai.ChatMessageRoleSystem,
			Content: "你是一个Quartz Cron表达式专家,我会向你进行描述，请根据我的描述生成表达式，并且只返回表达式，例如 0 30 7 1/1 * ?"},
		{Role: openai.ChatMessageRoleUser, Content: "每天早上八点准时推送"},
	}
	t.Log(clients.createChat(context.Background(), openai.GPT3Dot5Turbo, msgs))
}

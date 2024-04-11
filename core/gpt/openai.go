package gpt

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"sync"
	"wechat-gptbot/config"
	"wechat-gptbot/consts"
)

type openAiClient struct {
	sync.RWMutex
	cs map[string]*openai.Client
}

func (c *openAiClient) addClient(model string) *openai.Client {
	c.Lock()
	defer c.Unlock()
	if client, ok := c.cs[model]; ok {
		return client
	}
	var (
		clientConfig openai.ClientConfig
	)

	switch model {
	case openai.GPT3Dot5Turbo:
		clientConfig = openai.DefaultConfig(config.C.Gpt.TextConfig.AuthToken)
		clientConfig.BaseURL = config.C.Gpt.TextConfig.BaseURL
	case openai.CreateImageModelDallE3:
		clientConfig = openai.DefaultConfig(config.C.Gpt.ImageConfig.AuthToken)
		clientConfig.BaseURL = config.C.Gpt.ImageConfig.BaseURL
	default:
		clientConfig = openai.DefaultConfig(config.C.Gpt.TextConfig.AuthToken)
		clientConfig.BaseURL = config.C.Gpt.TextConfig.BaseURL
	}
	client := openai.NewClientWithConfig(clientConfig)
	c.cs[model] = client
	return client
}

func (c *openAiClient) getClient(model string) *openai.Client {
	c.RLock()
	if client, ok := c.cs[model]; ok {
		c.RUnlock()
		return client
	}
	c.RUnlock()
	return c.addClient(model)
}

// 发送聊天信息到Openai
func (c *openAiClient) createChat(ctx context.Context, model string, messages []openai.ChatCompletionMessage) string {
	resp, err := c.getClient(model).CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:           model,
		Messages:        messages,
		TopP:            1,
		Temperature:     0.5,
		PresencePenalty: 0,
	})
	if err != nil {
		logrus.Infof("openAIClient.CreateChatCompletion err=%+v\n", err)
		return consts.ErrTips
	}
	if len(resp.Choices) == 0 {
		logrus.Infof("resp is err, resp=%+v\n", resp)
		return consts.ErrTips
	}
	return resp.Choices[0].Message.Content
}

// 文生图模型 返回URL
func (c *openAiClient) createImage(ctx context.Context, model, prompt string) string {
	resp, err := c.getClient(model).CreateImage(ctx, openai.ImageRequest{
		Prompt:         prompt,
		N:              1,
		Quality:        openai.CreateImageQualityHD,
		Style:          openai.CreateImageStyleVivid,
		Size:           openai.CreateImageSize1024x1024,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		Model:          model,
		User:           "user",
	})
	if err != nil {
		logrus.Infof("openAIClient.CreateImage err=%+v\n", err)
		return ""
	}
	return resp.Data[0].URL
}

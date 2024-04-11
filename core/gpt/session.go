package gpt

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"reflect"
	"sync"
	"wechat-gptbot/config"
)

/*
* @Author: zouyx
* @Date:   2023/11/11 18:05
* @Package: 支持会话上下文管理 暂时只保留最近3次对话信息
 */
var messageCtx *session

// 存放用户上下文
type session struct {
	sync.RWMutex
	client *openAiClient           // 会话客户端
	ctx    map[string]*userMessage // 管理用户上下文
}

func InitSession() {
	clients := &openAiClient{}
	gptConfigValues := reflect.ValueOf(config.C.Gpt)
	numField := gptConfigValues.NumField()
	clients.cs = make(map[string]*openai.Client, numField)
	messageCtx = &session{
		ctx:    make(map[string]*userMessage),
		client: clients,
	}
}

// 用户级消息
type userMessage struct {
	mu   sync.Mutex                     // 加锁 防止上下文顺序紊乱  todo
	user string                         // 用户
	ctx  []openai.ChatCompletionMessage // 用户聊天的上下文 最多只保留6条记录，3组对话
}

// 新建一个用户级消息
func newUserMessage(user string, msg openai.ChatCompletionMessage) *userMessage {
	return &userMessage{
		user: user,
		ctx: []openai.ChatCompletionMessage{{
			Role:    openai.ChatMessageRoleSystem,
			Content: config.Prompt,
		}, msg},
		mu: sync.Mutex{},
	}
}

// 添加上下文
func (c *session) addContext(userName string, currentMessage openai.ChatCompletionMessage) {
	var (
		msg *userMessage
		ok  bool
	)
	if msg, ok = c.ctx[userName]; ok {
		// 直接追加到上下文中
		msg.ctx = append(msg.ctx, currentMessage)
		// 最多保存6条上下文
		if len(msg.ctx) > 6 {
			msg.ctx = msg.ctx[len(msg.ctx)-6:]
			// 将prompt 作为第一句传给机器人
			msg.ctx[0] = PromptMessage
		}
		return
	}
	// 当前没有上下文，新建一个用户级上下消息体
	c.ctx[userName] = newUserMessage(userName, currentMessage)
}

var PromptMessage openai.ChatCompletionMessage

// 构建上下文到消息体
func (c *session) buildMessage(userName, content string) []openai.ChatCompletionMessage {
	// 将当前对话加入上下文
	c.addContext(userName, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	})
	fmt.Println("=====" + userName + "=======")
	for i, ctx := range c.ctx[userName].ctx {
		fmt.Printf("%d     %s\n", i, ctx.Content)
	}
	fmt.Println("=====" + userName + "=======")
	return c.ctx[userName].ctx
}

func Chat(ctx context.Context, content string) string {
	// 默认不带上下文
	msgs := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: config.Prompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		},
	}
	sender := ctx.Value("sender").(string)
	if config.C.ContextStatus {
		msgs = messageCtx.buildMessage(sender, content)
	}
	// 发送消息
	reply := messageCtx.client.createChat(ctx, config.C.BaseModel, msgs)
	if config.C.ContextStatus {
		// 4. 把回复添加进上下文
		messageCtx.addContext(sender, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: reply,
		})
	}
	return reply
}

func CreateImage(ctx context.Context, prompt string) string {
	return messageCtx.client.createImage(ctx, openai.CreateImageModelDallE3, prompt)
}

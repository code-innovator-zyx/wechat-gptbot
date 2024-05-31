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

type Session interface {
	Chat(ctx context.Context, content string) string       // 对话
	CreateImage(ctx context.Context, prompt string) string // 生成图片，返回URL
}

// Session 存放用户上下文
type session struct {
	sync.Mutex                         // 用户的创建需要加锁
	client     *openAiClient           // 会话客户端
	ctx        map[string]*userMessage // 管理用户上下文
}

func NewSession() Session {
	clients := &openAiClient{}
	gptConfigValues := reflect.ValueOf(config.C.Gpt)
	numField := gptConfigValues.NumField()
	clients.cs = make(map[string]*openai.Client, numField)
	return &session{
		Mutex:  sync.Mutex{},
		ctx:    make(map[string]*userMessage),
		client: clients,
	}
}

// 获取用户
func (s *session) getUserContext(userName string) *userMessage {

	if msg, ok := s.ctx[userName]; ok {
		return msg
	}
	s.Lock()
	defer s.Unlock()
	// 双检加锁，防止加锁的过程中已经创建了用户
	if msg, ok := s.ctx[userName]; ok {
		return msg
	}
	msg := newUserMessage(userName)
	s.ctx[userName] = msg
	return msg
}

// 用户级消息
type userMessage struct {
	sync.Mutex                                // 加锁 防止上下文顺序紊乱 一个用户只能拿到响应后才能再次提问
	user       string                         // 用户
	ctx        []openai.ChatCompletionMessage // 用户聊天的上下文 最多只保留6条记录，3组对话
}

// 新建一个用户级消息
func newUserMessage(user string) *userMessage {
	return &userMessage{
		user: user,
		ctx: []openai.ChatCompletionMessage{{
			Role:    openai.ChatMessageRoleSystem,
			Content: config.Prompt,
		}},
		Mutex: sync.Mutex{},
	}
}

// 给用户追加上下文
func (um *userMessage) addContext(currentMessage openai.ChatCompletionMessage) {
	um.ctx = append(um.ctx, currentMessage)
	// 最多保存6条上下文
	if len(um.ctx) > 6 {
		um.ctx = um.ctx[len(um.ctx)-6:]
		// 将prompt 作为第一句传给机器人
		um.ctx[0] = PromptMessage
	}
}

// 构建上下文到消息体
func (um *userMessage) buildMessage(userName, content string) []openai.ChatCompletionMessage {
	// 将当前对话加入上下文
	um.addContext(openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	})
	fmt.Println("=====" + userName + "=======")
	for i, ctx := range um.ctx {
		fmt.Printf("%d     %s\n", i, ctx.Content)
	}
	fmt.Println("=====" + userName + "=======")
	return um.ctx
}

var PromptMessage openai.ChatCompletionMessage

func (s *session) Chat(ctx context.Context, content string) string {
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
	// 获取用户上下文
	um := s.getUserContext(sender)
	if config.C.ContextStatus {
		// 只有在用户开启上下文的时候，追加上下文需要加锁,得到回复追加上下文后才进行锁的释放
		um.Lock()
		defer um.Unlock()
		msgs = um.buildMessage(sender, content)
	}
	// 发送消息
	reply := s.client.createChat(ctx, config.C.GetBaseModel(), msgs)
	if config.C.ContextStatus {
		// 4. 把回复添加进上下文
		um.addContext(openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: reply,
		})
	}
	return reply
}

func (s *session) CreateImage(ctx context.Context, prompt string) string {
	return s.client.createImage(ctx, openai.CreateImageModelDallE3, prompt)
}

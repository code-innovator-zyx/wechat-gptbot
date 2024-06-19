package gpt

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"reflect"
	"sync"
	"wechat-gptbot/config"
	"wechat-gptbot/core/plugins"
)

/*
* @Author: zouyx
* @Date:   2023/11/11 18:05
* @Package: 支持会话上下文管理 暂时只保留最近3次对话信息
 */
const MaxSession = 6

type Session interface {
	Chat(ctx context.Context, content string) []string     // 对话
	CreateImage(ctx context.Context, prompt string) string // 生成图片，返回URL
}

// Session 存放用户上下文
type session struct {
	sync.RWMutex                                // 用户的创建需要加锁
	client         *openAiClient                // 会话客户端
	ctx            map[string]*userMessage      // 管理用户上下文
	prompt         openai.ChatCompletionMessage // 管理提示词
	pluginRegistry *plugins.PluginManger        // 插件注册器
}

func NewSession() Session {
	clients := &openAiClient{}
	gptConfigValues := reflect.ValueOf(config.C.Gpt)
	numField := gptConfigValues.NumField()
	clients.cs = make(map[string]*openai.Client, numField)
	registry := plugins.NewPluginRegistry()
	return &session{
		RWMutex:        sync.RWMutex{},
		ctx:            make(map[string]*userMessage),
		client:         clients,
		prompt:         initPrompt(),
		pluginRegistry: registry,
	}
}
func initPrompt() openai.ChatCompletionMessage {
	// 获取所有插件信息
	pluginsInfo := plugins.Manger.PluginPrompt()
	prompt := fmt.Sprintf(config.Prompt, pluginsInfo)
	return openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: prompt,
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
	msg := s.newUserMessage(userName)
	return msg
}

// Prompt 获取提示词  todo:将插件写入提示词
func (s *session) Prompt() openai.ChatCompletionMessage {
	return s.prompt
}

// 用户级消息
type userMessage struct {
	sync.Mutex                                // 加锁 防止上下文顺序紊乱 一个用户只能拿到响应后才能再次提问
	user       string                         // 用户
	ctx        []openai.ChatCompletionMessage // 用户聊天的上下文 最多只保留6条记录，3组对话
}

// 新建一个用户级消息
func (s *session) newUserMessage(user string) *userMessage {
	msg := &userMessage{
		user:  user,
		ctx:   []openai.ChatCompletionMessage{s.prompt},
		Mutex: sync.Mutex{},
	}
	s.ctx[user] = msg
	return msg
}

// 给用户追加上下文
func (um *userMessage) addContext(currentMessage, prompt openai.ChatCompletionMessage) {
	um.ctx = append(um.ctx, currentMessage)
	// 最多保存6条上下文
	if len(um.ctx) > MaxSession {
		um.ctx = um.ctx[len(um.ctx)-MaxSession:]
		// 将prompt 作为上下文第一条
		um.ctx[0] = prompt
	}
}

// 构建上下文到消息体
func (um *userMessage) buildMessage(userName string, currentMsg openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
	msgs := append(um.ctx, currentMsg)
	fmt.Println("=====" + userName + "=======")
	for i, ctx := range msgs {
		if i <= 0 {
			continue
		}
		fmt.Printf("%d     %s\n", i, ctx.Content)
	}
	fmt.Println("=====" + userName + "=======")
	return msgs
}

func (s *session) Chat(ctx context.Context, content string) []string {
	currentMsg := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	}
	// 默认不带上下文
	msgs := []openai.ChatCompletionMessage{
		s.Prompt(),
		currentMsg,
	}
	sender := ctx.Value("sender").(string)
	// 获取用户上下文
	um := s.getUserContext(sender)
	if config.C.ContextStatus {
		// 只有在用户开启上下文的时候，追加上下文需要加锁,得到回复追加上下文后才进行锁的释放
		um.Lock()
		defer um.Unlock()
		msgs = um.buildMessage(sender, currentMsg)
	}
	// 发送消息
	reply, err := s.client.createChat(ctx, config.C.GetBaseModel(), msgs)
	if nil != err {
		// 发送失败嘞
		return []string{err.Error()}
	}
	// 发送成功，可以讲请求和回复加入上下文
	if config.C.ContextStatus {
		// 如果请求成功才把问题回复都添加进上下文
		if resetMsg, ok := plugins.Manger.DoPlugin(reply); ok {
			reply = resetMsg
			goto RETURN
		}
		um.addContext(currentMsg, s.Prompt())
		um.addContext(openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: reply[0],
		}, s.Prompt())
	}
RETURN:
	return reply
}

func (s *session) CreateImage(ctx context.Context, prompt string) string {
	return s.client.createImage(ctx, openai.CreateImageModelDallE3, prompt)
}

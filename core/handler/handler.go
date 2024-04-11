package handler

import (
	"bytes"
	"context"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"net/url"
	"path"
	"strings"
	"wechat-gptbot/config"
	"wechat-gptbot/consts"
	"wechat-gptbot/core/gpt"
	"wechat-gptbot/utils"
)

func MessageHandler(msg *openwechat.Message) {
	// 判断是否需要我处理的对话
	if !checkMessageType(msg) {
		return
	}
	sender, err := msg.Sender()
	if nil != err {
		logrus.Error(err.Error())
	}

	if msg.IsSendByGroup() {
		sender, err = msg.SenderInGroup()
	}
	// 根据聊天类型使用不同的处理器去处理聊天逻辑  todo  不同用户的群聊信息和私信信息上下文隔离
	if msg.IsSendByGroup() {
		sender, err = msg.SenderInGroup()
		if nil != err {
			logrus.Error(err)
			return
		}
	}
	ctx := context.WithValue(context.TODO(), "sender", sender.NickName)
	childHandler := textReplyHandler
	if checkCreateImage(msg) {
		childHandler = imageReplyHandler
	}
	childHandler(ctx, msg)
}

// 判断是否是发给我的消息
func checkMessageType(msg *openwechat.Message) bool {
	if !msg.IsText() {
		// 目前只能处理文本对话
		return false
	}
	if !msg.IsSendByGroup() {
		// 私信消息
		return true
	}
	//  如果是艾特我的消息
	if msg.IsAt() {
		msg.Content = msg.Content[len("@年年"):]
		return true
	}
	// 如果包含了我的唤醒词
	if strings.HasPrefix(msg.Content, config.C.Gpt.TextConfig.TriggerPrefix) {
		msg.Content = strings.TrimPrefix(msg.Content, config.C.Gpt.TextConfig.TriggerPrefix)
		return true
	}

	if strings.HasPrefix(msg.Content, config.C.Gpt.ImageConfig.TriggerPrefix) {
		msg.Content = strings.TrimPrefix(msg.Content, config.C.Gpt.ImageConfig.TriggerPrefix)
		return true
	}

	return false
}

// 通过语义判断是否是文生图的需求
func checkCreateImage(msg *openwechat.Message) bool {
	if strings.HasPrefix(msg.Content, config.C.Gpt.ImageConfig.TriggerPrefix) {
		msg.Content = strings.TrimPrefix(msg.Content, config.C.Gpt.ImageConfig.TriggerPrefix)
		return true
	}
	return false
}

// 回复文本
func textReplyHandler(ctx context.Context, msg *openwechat.Message) {
	sender := ctx.Value("sender").(string)
	reply := gpt.Chat(ctx, utils.BuildPersonalMessage(sender, msg.Content))
	fmt.Printf("[text] Response: %s\n", reply) // 输出回复消息到日志
	_, err := msg.ReplyText(utils.BuildResponseMessage(sender, msg.Content, reply))
	if err != nil {
		logrus.Infof("msg.ReplyText Error: %+v", err)
	}
	msg.RevokeMsg()
}

// 回复图片
func imageReplyHandler(ctx context.Context, msg *openwechat.Message) {
	prompt := strings.TrimSpace(msg.Content)
	url := gpt.CreateImage(ctx, prompt)
	if url == "" {
		logrus.Infof("[image] Response: url 为空")
		msg.ReplyText(consts.ErrTips)
		return
	}
	logrus.Infof("[image] Response: url = %s", url)
	reader := bytes.Buffer{}
	err := utils.CompressImage(url, &reader)
	if err != nil {
		logrus.Infof("[image] downloadImage err, err=%+v", err)
		msg.ReplyText(consts.ErrTips)
		return
	}
	fu := msg.ReplyImage
	if checkFile(url) {
		fu = msg.ReplyFile
	}
	_, err = fu(&reader)
	if err != nil {
		logrus.Infof("msg.ReplyImage Error: %+v", err)
	}
}

func checkFile(uri string) bool {
	u, _ := url.Parse(uri)
	// 获取文件名
	name := path.Base(u.Path)
	return path.Ext(name) == ".webp"
}

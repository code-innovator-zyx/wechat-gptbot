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
	"wechat-gptbot/core/svc"
	"wechat-gptbot/utils"
)

var Context *svc.ServiceContext

type MessageMatchDispatcher struct {
	*openwechat.MessageMatchDispatcher
	ctx *svc.ServiceContext
}

func NewMessageMatchDispatcher() *MessageMatchDispatcher {
	dispatcher := openwechat.NewMessageMatchDispatcher()
	self := &MessageMatchDispatcher{
		dispatcher,
		svc.NewServiceContext(),
	}
	// 注册文本函数
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		needReply, isImage := checkMessageType(message)
		return needReply && !isImage
	}, self.text)

	// 注册图片函数
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		needReply, isImage := checkMessageType(message)
		return needReply && isImage
	}, self.image)
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		return message.IsTickled()
	}, self.trick)
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		return message.IsJoinGroup()
	}, self.joinGroup)
	dispatcher.SetAsync(true)
	return self
}
func (dispatcher *MessageMatchDispatcher) trick(message *openwechat.MessageContext) {
	message.ReplyText("别拍啦，小屁股都开花啦")
	return
}
func (dispatcher *MessageMatchDispatcher) joinGroup(message *openwechat.MessageContext) {
	message.ReplyText("欢迎欢迎，热烈欢迎")
	return
}

func (dispatcher *MessageMatchDispatcher) text(message *openwechat.MessageContext) {
	sender, err := message.Sender()
	if nil != err {
		logrus.Error(err.Error())
	}

	if message.IsSendByGroup() {
		sender, err = message.SenderInGroup()
	}
	for _, reply := range Context.Session.Chat(context.WithValue(context.TODO(), "sender", sender.NickName), utils.BuildPersonalMessage(sender.NickName, message.Content)) {
		fmt.Printf("[text] Response: %s\n", reply) // 输出回复消息到日志
		_, err = message.ReplyText(reply)
		if err != nil {
			logrus.Infof("msg.ReplyText Error: %+v", err)
		}
	}

}

func (dispatcher *MessageMatchDispatcher) image(message *openwechat.MessageContext) {
	message.Content = strings.TrimLeft(message.Content, config.C.Base.Gpt.ImageConfig.TriggerPrefix)
	sender, err := message.Sender()
	if nil != err {
		logrus.Error(err.Error())
	}

	if message.IsSendByGroup() {
		sender, err = message.SenderInGroup()
	}

	prompt := strings.TrimSpace(message.Content)
	uri := Context.Session.CreateImage(context.WithValue(context.TODO(), "sender", sender.NickName), prompt)
	if uri == "" {
		logrus.Infof("[image] Response: url 为空")
		message.ReplyText(consts.ErrTips)
		return
	}
	logrus.Infof("[image] Response: url = %s", uri)
	reader := bytes.Buffer{}
	err = utils.CompressImage(uri, &reader)
	if err != nil {
		logrus.Infof("[image] downloadImage err, err=%+v", err)
		message.ReplyText(consts.ErrTips)
		return
	}
	fu := message.ReplyImage
	if checkFile(uri) {
		fu = message.ReplyFile
	}
	_, err = fu(&reader)
	if err != nil {
		logrus.Infof("msg.ReplyImage Error: %+v", err)
	}
}

// 判断是否是发给我的消息
func checkMessageType(msg *openwechat.Message) (needReply bool, isImage bool) {
	// 如果包含了我的唤醒词
	msg.Content = strings.TrimLeft(msg.Content, " ")
	sender, err := msg.Sender()
	if nil != err {
		logrus.Error(err.Error())
	}
	if !msg.IsText() {
		return false, false
	}
	if msg.IsSendBySelf() {
		return false, false
	}
	if !msg.IsSendByGroup() {
		// 私信消息
		// 私信消息不要管公众号消息
		if sender.IsMP() {
			return false, false
		}
		return true, checkCreateImage(msg)
	}
	//  如果是艾特我的消息
	if msg.IsAt() {
		prefix := fmt.Sprintf("@%s", msg.Owner().NickName)
		if strings.HasPrefix(msg.Content, prefix) {
			msg.Content = msg.Content[len(prefix):]
		}
		return true, checkCreateImage(msg)
	}

	if strings.HasPrefix(msg.Content, config.C.Base.Gpt.TextConfig.TriggerPrefix) {
		msg.Content = strings.TrimLeft(msg.Content, config.C.Base.Gpt.TextConfig.TriggerPrefix)
		return true, false
	}

	if strings.HasPrefix(msg.Content, config.C.Base.Gpt.ImageConfig.TriggerPrefix) {
		msg.Content = strings.TrimLeft(msg.Content, config.C.Base.Gpt.ImageConfig.TriggerPrefix)
		return true, true
	}

	return false, false
}

// 通过语义判断是否是文生图的需求
func checkCreateImage(msg *openwechat.Message) bool {
	msg.Content = strings.TrimPrefix(msg.Content, "\u2005")
	if strings.HasPrefix(msg.Content, config.C.Base.Gpt.ImageConfig.TriggerPrefix) {
		return true
	}
	return false
}

func checkFile(uri string) bool {
	u, _ := url.Parse(uri)
	// 获取文件名
	name := path.Base(u.Path)
	return path.Ext(name) == ".webp"
}

package main

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"wechat-gptbot/core"
	"wechat-gptbot/core/handler"
)

/*
* @Author: zouyx
* @Email: zouyx@knowsec.com
* @Date:   2024/4/8 10:46
* @Package:
 */
func main() {
	// 初始化核心配置
	core.Initialize()
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式
	// 定义消息处理函数
	bot.MessageHandler = handler.MessageHandler
	bot.UUIDCallback = consoleQrCode // 注册登陆二维码回调
	// 登录回调
	bot.SyncCheckCallback = nil
	reloadStorage := openwechat.NewFileHotReloadStorage("token.json")
	if err := bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption()); nil != err {
		panic(err)
	}
	// 获取当前登录的用户
	self, err := bot.GetCurrentUser()
	if nil != err {
		panic(err)
	}
	logrus.Infof("login success %+v", *self.User)
	bot.Block()
}

func consoleQrCode(uuid string) {
	println("访问下面网址扫描二维码登录")
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Medium)
	fmt.Println(q.ToSmallString(false))
}
